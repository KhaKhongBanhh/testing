package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"movie-service/internal/module/movie/entity"

	rediscache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// ---------- fakes ----------

type fakeCache struct{}

func (fakeCache) Get(context.Context, string, any) error { return rediscache.ErrCacheMiss }
func (fakeCache) Set(context.Context, string, any, time.Duration) error { return nil }
func (fakeCache) Delete(context.Context, string) error                  { return nil }

// fakeRedis only supports what caching.DeleteKeys needs (Scan/Del).
type fakeRedis struct{ redis.UniversalClient }

func (fakeRedis) Scan(ctx context.Context, _ uint64, _ string, _ int64) *redis.ScanCmd {
	c := redis.NewScanCmd(ctx, nil)
	c.SetVal([]string{}, 0)
	return c
}
func (fakeRedis) Del(ctx context.Context, _ ...string) *redis.IntCmd { return redis.NewIntCmd(ctx) }

type fakeRepo struct {
	movie       *entity.Movie
	getErr      error
	createErr   error
	updateErr   error
	deleteErr   error
	created     int
	updated     *entity.Movie
	updatedGen  []string
	deleted     string
	total       int
	many        []*entity.Movie
	gotLimit    int
	gotOffset   int
}

func (f *fakeRepo) GetByID(context.Context, string) (*entity.Movie, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	cp := *f.movie
	return &cp, nil
}
func (f *fakeRepo) GetMany(_ context.Context, limit, offset int, _ string, _ string) ([]*entity.Movie, error) {
	f.gotLimit, f.gotOffset = limit, offset
	return f.many, nil
}
func (f *fakeRepo) GetTotalCount(context.Context, string, string) (int, error) { return f.total, nil }
func (f *fakeRepo) GetMovieStats(context.Context) ([]*entity.MovieStat, error) {
	return []*entity.MovieStat{{Status: "SHOWING", Count: 2}}, nil
}
func (f *fakeRepo) GetGenres(context.Context) ([]*entity.Genre, error) {
	return []*entity.Genre{{Id: "g1", Name: "Action"}}, nil
}
func (f *fakeRepo) Create(_ context.Context, _ *entity.Movie, _ []string) error {
	f.created++
	return f.createErr
}
func (f *fakeRepo) Update(_ context.Context, m *entity.Movie, g []string) error {
	f.updated, f.updatedGen = m, g
	return f.updateErr
}
func (f *fakeRepo) Delete(_ context.Context, id string) error {
	f.deleted = id
	return f.deleteErr
}

func newBiz(r *fakeRepo) *business {
	return &business{repository: r, cache: fakeCache{}, roCache: fakeCache{}, redisClient: fakeRedis{}}
}

var ctx = context.Background()

// ---------- tests ----------

func TestCreateMovie(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-01a", func(t *testing.T) { // nil movie
		if err := newBiz(&fakeRepo{}).CreateMovie(ctx, nil, nil); !errors.Is(err, ErrInvalidMovieData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-01b", func(t *testing.T) { // empty title
		r := &fakeRepo{}
		err := newBiz(r).CreateMovie(ctx, &entity.Movie{Title: "", Duration: 100}, nil)
		if !errors.Is(err, ErrInvalidMovieData) || r.created != 0 {
			t.Fatalf("err=%v created=%d", err, r.created)
		}
	})
	t.Run("TC-MOVIE-BIZ-01c", func(t *testing.T) { // duration 0
		r := &fakeRepo{}
		err := newBiz(r).CreateMovie(ctx, &entity.Movie{Title: "A", Duration: 0}, nil)
		if !errors.Is(err, ErrInvalidMovieData) || r.created != 0 {
			t.Fatalf("err=%v created=%d", err, r.created)
		}
	})
	t.Run("TC-MOVIE-BIZ-01d", func(t *testing.T) { // empty status gets default
		r := &fakeRepo{}
		m := &entity.Movie{Title: "A", Duration: 100}
		if err := newBiz(r).CreateMovie(ctx, m, []string{"g1"}); err != nil {
			t.Fatal(err)
		}
		if m.Status != entity.MovieStatusUpcoming || r.created != 1 {
			t.Fatalf("status=%s created=%d", m.Status, r.created)
		}
	})
	t.Run("TC-MOVIE-BIZ-01e", func(t *testing.T) { // explicit status preserved
		m := &entity.Movie{Title: "A", Duration: 100, Status: entity.MovieStatusShowing}
		if err := newBiz(&fakeRepo{}).CreateMovie(ctx, m, nil); err != nil {
			t.Fatal(err)
		}
		if m.Status != entity.MovieStatusShowing {
			t.Fatalf("status=%s", m.Status)
		}
	})
	t.Run("TC-MOVIE-BIZ-01f", func(t *testing.T) { // repository error wrapped
		boom := errors.New("db down")
		err := newBiz(&fakeRepo{createErr: boom}).CreateMovie(ctx, &entity.Movie{Title: "A", Duration: 1}, nil)
		if !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestUpdateMovie(t *testing.T) {
	existing := func(s entity.MovieStatus) *entity.Movie {
		return &entity.Movie{Id: "m1", Title: "Old", Duration: 100, Status: s}
	}
	t.Run("TC-MOVIE-BIZ-02a", func(t *testing.T) { // nil movie
		if err := newBiz(&fakeRepo{}).UpdateMovie(ctx, nil, nil); !errors.Is(err, ErrInvalidMovieData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-02b", func(t *testing.T) { // empty id
		err := newBiz(&fakeRepo{}).UpdateMovie(ctx, &entity.Movie{Title: "A", Duration: 1}, nil)
		if !errors.Is(err, ErrInvalidMovieData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-02c", func(t *testing.T) { // invalid data
		err := newBiz(&fakeRepo{}).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "", Duration: 1}, nil)
		if !errors.Is(err, ErrInvalidMovieData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-02d", func(t *testing.T) { // movie not found
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "A", Duration: 1}, nil)
		if !errors.Is(err, ErrMovieNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-02e", func(t *testing.T) { // same status -> allowed
		r := &fakeRepo{movie: existing(entity.MovieStatusUpcoming)}
		err := newBiz(r).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "New", Duration: 100, Status: entity.MovieStatusUpcoming}, []string{"g1"})
		if err != nil || r.updated == nil || r.updated.Title != "New" {
			t.Fatalf("err=%v updated=%v", err, r.updated)
		}
	})
	t.Run("TC-MOVIE-BIZ-02f", func(t *testing.T) { // valid transition
		r := &fakeRepo{movie: existing(entity.MovieStatusUpcoming)}
		err := newBiz(r).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "A", Duration: 100, Status: entity.MovieStatusShowing}, nil)
		if err != nil || r.updated == nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-02g", func(t *testing.T) { // invalid transition SHOWING -> UPCOMING
		r := &fakeRepo{movie: existing(entity.MovieStatusShowing)}
		err := newBiz(r).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "A", Duration: 100, Status: entity.MovieStatusUpcoming}, nil)
		if !errors.Is(err, ErrInvalidStatusTransition) || r.updated != nil {
			t.Fatalf("err=%v updated=%v", err, r.updated)
		}
	})
	t.Run("TC-MOVIE-BIZ-02h", func(t *testing.T) { // omitted status must keep the current status (expected behaviour)
		r := &fakeRepo{movie: existing(entity.MovieStatusShowing)}
		err := newBiz(r).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "A", Duration: 100, Status: ""}, nil)
		if err != nil {
			t.Fatalf("updating without a status must not fail, got: %v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-02i", func(t *testing.T) { // repository error
		boom := errors.New("db down")
		r := &fakeRepo{movie: existing(entity.MovieStatusUpcoming), updateErr: boom}
		err := newBiz(r).UpdateMovie(ctx, &entity.Movie{Id: "m1", Title: "A", Duration: 1, Status: entity.MovieStatusUpcoming}, nil)
		if !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestDeleteMovie(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-03a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).DeleteMovie(ctx, ""); err == nil {
			t.Fatal("want error for empty id")
		}
	})
	t.Run("TC-MOVIE-BIZ-03b", func(t *testing.T) {
		r := &fakeRepo{getErr: sql.ErrNoRows}
		if err := newBiz(r).DeleteMovie(ctx, "m1"); !errors.Is(err, ErrMovieNotFound) || r.deleted != "" {
			t.Fatalf("err=%v deleted=%q", err, r.deleted)
		}
	})
	t.Run("TC-MOVIE-BIZ-03c", func(t *testing.T) {
		boom := errors.New("fk violation")
		r := &fakeRepo{movie: &entity.Movie{Id: "m1"}, deleteErr: boom}
		if err := newBiz(r).DeleteMovie(ctx, "m1"); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-03d", func(t *testing.T) {
		r := &fakeRepo{movie: &entity.Movie{Id: "m1"}}
		if err := newBiz(r).DeleteMovie(ctx, "m1"); err != nil || r.deleted != "m1" {
			t.Fatalf("err=%v deleted=%q", err, r.deleted)
		}
	})
}

func TestUpdateMovieStatus(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-04a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).UpdateMovieStatus(ctx, "", entity.MovieStatusShowing); err == nil {
			t.Fatal("want error")
		}
	})
	t.Run("TC-MOVIE-BIZ-04b", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).UpdateMovieStatus(ctx, "m1", entity.MovieStatusShowing)
		if !errors.Is(err, ErrMovieNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-04c", func(t *testing.T) { // UPCOMING -> SHOWING
		r := &fakeRepo{movie: &entity.Movie{Id: "m1", Status: entity.MovieStatusUpcoming}}
		if err := newBiz(r).UpdateMovieStatus(ctx, "m1", entity.MovieStatusShowing); err != nil {
			t.Fatal(err)
		}
		if r.updated == nil || r.updated.Status != entity.MovieStatusShowing || r.updatedGen != nil {
			t.Fatalf("updated=%+v genres=%v", r.updated, r.updatedGen)
		}
	})
	t.Run("TC-MOVIE-BIZ-04d", func(t *testing.T) { // SHOWING -> UPCOMING rejected
		r := &fakeRepo{movie: &entity.Movie{Id: "m1", Status: entity.MovieStatusShowing}}
		err := newBiz(r).UpdateMovieStatus(ctx, "m1", entity.MovieStatusUpcoming)
		if !errors.Is(err, ErrInvalidStatusTransition) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-04e", func(t *testing.T) { // UPCOMING -> ENDED rejected
		r := &fakeRepo{movie: &entity.Movie{Id: "m1", Status: entity.MovieStatusUpcoming}}
		if err := newBiz(r).UpdateMovieStatus(ctx, "m1", entity.MovieStatusEnded); !errors.Is(err, ErrInvalidStatusTransition) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-04f", func(t *testing.T) { // ENDED is terminal
		r := &fakeRepo{movie: &entity.Movie{Id: "m1", Status: entity.MovieStatusEnded}}
		if err := newBiz(r).UpdateMovieStatus(ctx, "m1", entity.MovieStatusShowing); !errors.Is(err, ErrInvalidStatusTransition) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestGetMovieById(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-05a", func(t *testing.T) {
		if _, err := newBiz(&fakeRepo{}).GetMovieById(ctx, ""); !errors.Is(err, ErrInvalidMovieData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-05b", func(t *testing.T) {
		if _, err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).GetMovieById(ctx, "x"); !errors.Is(err, ErrMovieNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-05c", func(t *testing.T) {
		m, err := newBiz(&fakeRepo{movie: &entity.Movie{Id: "m1", Title: "A"}}).GetMovieById(ctx, "m1")
		if err != nil || m.Id != "m1" {
			t.Fatalf("m=%v err=%v", m, err)
		}
	})
	t.Run("TC-MOVIE-BIZ-05d", func(t *testing.T) { // other repository error is wrapped, not mapped to NotFound
		boom := errors.New("db down")
		_, err := newBiz(&fakeRepo{getErr: boom}).GetMovieById(ctx, "m1")
		if !errors.Is(err, boom) || errors.Is(err, ErrMovieNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestGetMovies(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-06a", func(t *testing.T) { // page/size normalised
		r := &fakeRepo{total: 3}
		if _, _, err := newBiz(r).GetMovies(ctx, 0, 0, "", ""); err != nil {
			t.Fatal(err)
		}
		if r.gotLimit != 10 || r.gotOffset != 0 {
			t.Fatalf("limit=%d offset=%d, want 10/0", r.gotLimit, r.gotOffset)
		}
	})
	t.Run("TC-MOVIE-BIZ-06b", func(t *testing.T) { // size capped at 100
		r := &fakeRepo{}
		if _, _, err := newBiz(r).GetMovies(ctx, 1, 500, "", ""); err != nil {
			t.Fatal(err)
		}
		if r.gotLimit != 100 {
			t.Fatalf("limit=%d, want 100", r.gotLimit)
		}
	})
	t.Run("TC-MOVIE-BIZ-06c", func(t *testing.T) { // offset computed
		r := &fakeRepo{}
		if _, _, err := newBiz(r).GetMovies(ctx, 3, 10, "", ""); err != nil {
			t.Fatal(err)
		}
		if r.gotOffset != 20 {
			t.Fatalf("offset=%d, want 20", r.gotOffset)
		}
	})
	t.Run("TC-MOVIE-BIZ-06d", func(t *testing.T) { // returns movies + total
		r := &fakeRepo{total: 7, many: []*entity.Movie{{Id: "a"}}}
		movies, total, err := newBiz(r).GetMovies(ctx, 1, 10, "a", "SHOWING")
		if err != nil || total != 7 || len(movies) != 1 {
			t.Fatalf("movies=%v total=%d err=%v", movies, total, err)
		}
	})
}

func TestValidateMovieForShowtime(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-07a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).ValidateMovieForShowtime(ctx, ""); !errors.Is(err, ErrInvalidMovieData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-07b", func(t *testing.T) {
		if err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).ValidateMovieForShowtime(ctx, "m1"); !errors.Is(err, ErrMovieNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-MOVIE-BIZ-07c", func(t *testing.T) {
		for _, s := range []entity.MovieStatus{entity.MovieStatusUpcoming, entity.MovieStatusEnded} {
			r := &fakeRepo{movie: &entity.Movie{Id: "m1", Status: s}}
			if err := newBiz(r).ValidateMovieForShowtime(ctx, "m1"); !errors.Is(err, ErrMovieNotShowing) {
				t.Fatalf("status %s: err=%v", s, err)
			}
		}
	})
	t.Run("TC-MOVIE-BIZ-07d", func(t *testing.T) {
		r := &fakeRepo{movie: &entity.Movie{Id: "m1", Status: entity.MovieStatusShowing}}
		if err := newBiz(r).ValidateMovieForShowtime(ctx, "m1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestGetMovieStatsAndGenres(t *testing.T) {
	t.Run("TC-MOVIE-BIZ-08a", func(t *testing.T) {
		stats, err := newBiz(&fakeRepo{}).GetMovieStats(ctx)
		if err != nil || len(stats) != 1 || stats[0].Count != 2 {
			t.Fatalf("stats=%v err=%v", stats, err)
		}
	})
	t.Run("TC-MOVIE-BIZ-08b", func(t *testing.T) {
		g, err := newBiz(&fakeRepo{}).GetGenres(ctx)
		if err != nil || len(g) != 1 || g[0].Name != "Action" {
			t.Fatalf("genres=%v err=%v", g, err)
		}
	})
}

