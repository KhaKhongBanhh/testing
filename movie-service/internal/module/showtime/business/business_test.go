package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	movieBusiness "movie-service/internal/module/movie/business"
	roomBusiness "movie-service/internal/module/room/business"
	"movie-service/internal/module/showtime/entity"

	rediscache "github.com/go-redis/cache/v9"
)

type fakeCache struct{}

func (fakeCache) Get(context.Context, string, any) error { return rediscache.ErrCacheMiss }
func (fakeCache) Set(context.Context, string, any, time.Duration) error { return nil }
func (fakeCache) Delete(context.Context, string) error                  { return nil }

type fakeMovieBiz struct {
	movieBusiness.MovieBiz
	err error
}

func (f fakeMovieBiz) ValidateMovieForShowtime(context.Context, string) error { return f.err }

type fakeRoomBiz struct {
	roomBusiness.RoomBiz
	err error
}

func (f fakeRoomBiz) ValidateRoomForShowtime(context.Context, string) error { return f.err }

type fakeRepo struct {
	st          *entity.Showtime
	getErr      error
	conflict    bool
	conflictErr error
	createErr   error
	updateErr   error
	deleteErr   error
	created     *entity.Showtime
	updated     *entity.Showtime
	deleted     string
	gotExcl     string
	conflictN   int
}

func (f *fakeRepo) GetByID(context.Context, string) (*entity.Showtime, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	cp := *f.st
	return &cp, nil
}
func (f *fakeRepo) GetByIds(context.Context, []string) ([]*entity.Showtime, error) { return nil, nil }
func (f *fakeRepo) GetMany(context.Context, int, int, string, string, string, entity.ShowtimeFormat, entity.ShowtimeStatus, *time.Time, *time.Time, bool) ([]*entity.Showtime, error) {
	return []*entity.Showtime{{Id: "a"}}, nil
}
func (f *fakeRepo) GetTotalCount(context.Context, string, string, string, entity.ShowtimeFormat, entity.ShowtimeStatus, *time.Time, *time.Time, bool) (int, error) {
	return 1, nil
}
func (f *fakeRepo) GetByMovie(context.Context, string) ([]*entity.Showtime, error) { return nil, nil }
func (f *fakeRepo) GetUpcoming(context.Context, int) ([]*entity.Showtime, error)    { return nil, nil }
func (f *fakeRepo) Create(_ context.Context, s *entity.Showtime) error {
	f.created = s
	return f.createErr
}
func (f *fakeRepo) Update(_ context.Context, s *entity.Showtime) error {
	f.updated = s
	return f.updateErr
}
func (f *fakeRepo) Delete(_ context.Context, id string) error { f.deleted = id; return f.deleteErr }
func (f *fakeRepo) CheckConflict(_ context.Context, _ string, _, _ time.Time, excl string) (bool, error) {
	f.gotExcl = excl
	f.conflictN++
	return f.conflict, f.conflictErr
}

func newBiz(r *fakeRepo, movieErr, roomErr error) *business {
	return &business{
		repository: r, movieBiz: fakeMovieBiz{err: movieErr}, roomBiz: fakeRoomBiz{err: roomErr},
		cache: fakeCache{}, roCache: fakeCache{},
	}
}

var ctx = context.Background()

func future(h int) time.Time { return time.Now().Add(time.Duration(h) * time.Hour).Truncate(time.Minute) }

func newSt() *entity.Showtime {
	return &entity.Showtime{Id: "st1", MovieId: "m1", RoomId: "r1", BasePrice: 80000,
		StartTime: future(48), EndTime: future(50), Status: entity.ShowtimeStatusScheduled}
}

func TestCreateShowtime(t *testing.T) {
	t.Run("TC-SHOWTIME-BIZ-01a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}, nil, nil).CreateShowtime(ctx, nil); !errors.Is(err, ErrInvalidShowtimeData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01b", func(t *testing.T) { // end before start
		s := newSt()
		s.EndTime = s.StartTime.Add(-time.Hour)
		if err := newBiz(&fakeRepo{}, nil, nil).CreateShowtime(ctx, s); !errors.Is(err, ErrInvalidShowtimeData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01c", func(t *testing.T) { // start in the past
		s := newSt()
		s.StartTime, s.EndTime = time.Now().Add(-48*time.Hour), time.Now().Add(-46*time.Hour)
		r := &fakeRepo{}
		if err := newBiz(r, nil, nil).CreateShowtime(ctx, s); !errors.Is(err, ErrShowtimeInPast) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01d", func(t *testing.T) { // movie not SHOWING
		r := &fakeRepo{}
		err := newBiz(r, movieBusiness.ErrMovieNotShowing, nil).CreateShowtime(ctx, newSt())
		if !errors.Is(err, ErrMovieNotShowing) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01e", func(t *testing.T) { // movie not found is passed through
		err := newBiz(&fakeRepo{}, movieBusiness.ErrMovieNotFound, nil).CreateShowtime(ctx, newSt())
		if !errors.Is(err, movieBusiness.ErrMovieNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01f", func(t *testing.T) { // room INACTIVE
		r := &fakeRepo{}
		err := newBiz(r, nil, roomBusiness.ErrRoomNotActive).CreateShowtime(ctx, newSt())
		if !errors.Is(err, ErrRoomNotActive) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01g", func(t *testing.T) { // room not found is passed through
		err := newBiz(&fakeRepo{}, nil, roomBusiness.ErrRoomNotFound).CreateShowtime(ctx, newSt())
		if !errors.Is(err, roomBusiness.ErrRoomNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01h", func(t *testing.T) { // time conflict
		r := &fakeRepo{conflict: true}
		err := newBiz(r, nil, nil).CreateShowtime(ctx, newSt())
		if !errors.Is(err, ErrTimeConflict) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01i", func(t *testing.T) { // conflict check error
		boom := errors.New("db")
		if err := newBiz(&fakeRepo{conflictErr: boom}, nil, nil).CreateShowtime(ctx, newSt()); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01j", func(t *testing.T) { // create error
		boom := errors.New("db")
		if err := newBiz(&fakeRepo{createErr: boom}, nil, nil).CreateShowtime(ctx, newSt()); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-01k", func(t *testing.T) { // success; conflict check excludes nothing
		r := &fakeRepo{}
		if err := newBiz(r, nil, nil).CreateShowtime(ctx, newSt()); err != nil {
			t.Fatal(err)
		}
		if r.created == nil || r.gotExcl != "" || r.conflictN != 1 {
			t.Fatalf("created=%v excl=%q n=%d", r.created, r.gotExcl, r.conflictN)
		}
	})
}

func TestUpdateShowtime(t *testing.T) {
	price := func(p float64) *float64 { return &p }
	str := func(s string) *string { return &s }
	t.Run("TC-SHOWTIME-BIZ-02a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}, nil, nil).UpdateShowtime(ctx, "", &entity.UpdateShowtimeRequest{}); !errors.Is(err, ErrInvalidShowtimeData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02b", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}, nil, nil).UpdateShowtime(ctx, "st1", nil); !errors.Is(err, ErrInvalidShowtimeData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02c", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{})
		if !errors.Is(err, ErrShowtimeNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02d", func(t *testing.T) { // price only: no conflict check
		r := &fakeRepo{st: newSt()}
		if err := newBiz(r, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{BasePrice: price(90000)}); err != nil {
			t.Fatal(err)
		}
		if r.updated.BasePrice != 90000 || r.conflictN != 0 {
			t.Fatalf("updated=%+v conflictN=%d", r.updated, r.conflictN)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02e", func(t *testing.T) { // new movie not SHOWING
		r := &fakeRepo{st: newSt()}
		err := newBiz(r, movieBusiness.ErrMovieNotShowing, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{MovieId: str("m2")})
		if !errors.Is(err, ErrMovieNotShowing) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02f", func(t *testing.T) { // new room INACTIVE
		r := &fakeRepo{st: newSt()}
		err := newBiz(r, nil, roomBusiness.ErrRoomNotActive).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{RoomId: str("r2")})
		if !errors.Is(err, ErrRoomNotActive) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02g", func(t *testing.T) { // end before start
		r := &fakeRepo{st: newSt()}
		e := future(1)
		err := newBiz(r, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{EndTime: &e})
		if !errors.Is(err, ErrInvalidShowtimeData) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02h", func(t *testing.T) { // move a SCHEDULED showtime into the past
		r := &fakeRepo{st: newSt()}
		s, e := time.Now().Add(-5*time.Hour), time.Now().Add(-3*time.Hour)
		err := newBiz(r, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{StartTime: &s, EndTime: &e})
		if !errors.Is(err, ErrShowtimeInPast) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02i", func(t *testing.T) { // reschedule conflicts; excludes itself
		r := &fakeRepo{st: newSt(), conflict: true}
		s, e := future(72), future(74)
		err := newBiz(r, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{StartTime: &s, EndTime: &e})
		if !errors.Is(err, ErrTimeConflict) || r.gotExcl != "st1" || r.updated != nil {
			t.Fatalf("err=%v excl=%q", err, r.gotExcl)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02j", func(t *testing.T) { // reschedule OK: times are rounded to 30 min
		r := &fakeRepo{st: newSt()}
		s := time.Date(2030, 1, 1, 10, 17, 0, 0, time.UTC)
		e := time.Date(2030, 1, 1, 12, 47, 0, 0, time.UTC)
		if err := newBiz(r, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{StartTime: &s, EndTime: &e}); err != nil {
			t.Fatal(err)
		}
		if r.updated.StartTime.Minute() != 0 || r.updated.EndTime.Minute() != 30 {
			t.Fatalf("start=%v end=%v", r.updated.StartTime, r.updated.EndTime)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-02k", func(t *testing.T) { // repository error
		boom := errors.New("db")
		r := &fakeRepo{st: newSt(), updateErr: boom}
		if err := newBiz(r, nil, nil).UpdateShowtime(ctx, "st1", &entity.UpdateShowtimeRequest{BasePrice: price(1)}); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestDeleteShowtime(t *testing.T) {
	t.Run("TC-SHOWTIME-BIZ-03a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}, nil, nil).DeleteShowtime(ctx, ""); !errors.Is(err, ErrInvalidShowtimeData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-03b", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}, nil, nil).DeleteShowtime(ctx, "st1")
		if !errors.Is(err, ErrShowtimeNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-03c", func(t *testing.T) {
		r := &fakeRepo{st: newSt()}
		if err := newBiz(r, nil, nil).DeleteShowtime(ctx, "st1"); err != nil || r.deleted != "st1" {
			t.Fatalf("err=%v deleted=%q", err, r.deleted)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-03d", func(t *testing.T) {
		boom := errors.New("fk booking")
		if err := newBiz(&fakeRepo{st: newSt(), deleteErr: boom}, nil, nil).DeleteShowtime(ctx, "st1"); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestUpdateShowtimeStatus(t *testing.T) {
	t.Run("TC-SHOWTIME-BIZ-04a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}, nil, nil).UpdateShowtimeStatus(ctx, "", entity.ShowtimeStatusOngoing); !errors.Is(err, ErrInvalidShowtimeData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-04b", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}, nil, nil).UpdateShowtimeStatus(ctx, "st1", entity.ShowtimeStatusOngoing)
		if !errors.Is(err, ErrShowtimeNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-04c", func(t *testing.T) { // SCHEDULED -> ONGOING
		r := &fakeRepo{st: newSt()}
		if err := newBiz(r, nil, nil).UpdateShowtimeStatus(ctx, "st1", entity.ShowtimeStatusOngoing); err != nil {
			t.Fatal(err)
		}
		if r.updated.Status != entity.ShowtimeStatusOngoing {
			t.Fatalf("status=%s", r.updated.Status)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-04d", func(t *testing.T) { // COMPLETED is final: reviving it must be rejected (expected behaviour)
		st := newSt()
		st.Status = entity.ShowtimeStatusCompleted
		r := &fakeRepo{st: st}
		err := newBiz(r, nil, nil).UpdateShowtimeStatus(ctx, "st1", entity.ShowtimeStatusScheduled)
		if err == nil {
			t.Fatalf("COMPLETED -> SCHEDULED was accepted")
		}
	})
	t.Run("TC-SHOWTIME-BIZ-04e", func(t *testing.T) { // unknown status must be rejected (expected behaviour)
		r := &fakeRepo{st: newSt()}
		if err := newBiz(r, nil, nil).UpdateShowtimeStatus(ctx, "st1", entity.ShowtimeStatus("FOO")); err == nil {
			t.Fatalf("status FOO was accepted and persisted")
		}
	})
}

func TestGetShowtimesAndConflict(t *testing.T) {
	t.Run("TC-SHOWTIME-BIZ-05a", func(t *testing.T) {
		for _, pr := range [][2]int{{0, 10}, {1, 0}} {
			if _, _, err := newBiz(&fakeRepo{}, nil, nil).GetShowtimes(ctx, pr[0], pr[1], "", "", "", "", "", nil, nil, false); !errors.Is(err, ErrInvalidShowtimeData) {
				t.Fatalf("page=%d size=%d err=%v", pr[0], pr[1], err)
			}
		}
	})
	t.Run("TC-SHOWTIME-BIZ-05b", func(t *testing.T) {
		list, total, err := newBiz(&fakeRepo{}, nil, nil).GetShowtimes(ctx, 1, 10, "", "", "", "", "", nil, nil, false)
		if err != nil || len(list) != 1 || total != 1 {
			t.Fatalf("list=%v total=%d err=%v", list, total, err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-05c", func(t *testing.T) {
		got, err := newBiz(&fakeRepo{conflict: true}, nil, nil).CheckTimeConflict(ctx, "r1", future(1), future(2), "x")
		if err != nil || !got {
			t.Fatalf("got=%v err=%v", got, err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-05d", func(t *testing.T) {
		if _, err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}, nil, nil).GetShowtimeById(ctx, "x"); !errors.Is(err, ErrShowtimeNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SHOWTIME-BIZ-05e", func(t *testing.T) {
		got, err := newBiz(&fakeRepo{st: newSt()}, nil, nil).GetShowtimesByIds(ctx, nil)
		if err != nil || len(got) != 0 {
			t.Fatalf("got=%v err=%v", got, err)
		}
	})
}

