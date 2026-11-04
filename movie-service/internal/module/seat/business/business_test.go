package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"movie-service/internal/module/seat/entity"

	rediscache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type fakeCache struct{}

func (fakeCache) Get(context.Context, string, any) error { return rediscache.ErrCacheMiss }
func (fakeCache) Set(context.Context, string, any, time.Duration) error { return nil }
func (fakeCache) Delete(context.Context, string) error                  { return nil }

type fakeRedis struct{ redis.UniversalClient }

func (fakeRedis) Scan(ctx context.Context, _ uint64, _ string, _ int64) *redis.ScanCmd {
	c := redis.NewScanCmd(ctx, nil)
	c.SetVal([]string{}, 0)
	return c
}
func (fakeRedis) Del(ctx context.Context, _ ...string) *redis.IntCmd { return redis.NewIntCmd(ctx) }

type fakeRepo struct {
	seat      *entity.Seat
	getErr    error
	exists    bool
	existsErr error
	createErr error
	updateErr error
	deleteErr error
	created   *entity.Seat
	updated   *entity.Seat
	deleted   string
	gotRow    string
	gotNum    string
	gotExcl   string
	gotOffset int
}

func (f *fakeRepo) GetByID(context.Context, string) (*entity.Seat, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	cp := *f.seat
	return &cp, nil
}
func (f *fakeRepo) GetByIDs(context.Context, []string) ([]*entity.Seat, error) {
	return []*entity.Seat{{Id: "a"}}, f.getErr
}
func (f *fakeRepo) GetMany(_ context.Context, _, offset int, _, _, _ string, _ entity.SeatType, _ entity.SeatStatus) ([]*entity.Seat, error) {
	f.gotOffset = offset
	return []*entity.Seat{{Id: "a"}}, nil
}
func (f *fakeRepo) GetTotalCount(context.Context, string, string, string, entity.SeatType, entity.SeatStatus) (int, error) {
	return 1, nil
}
func (f *fakeRepo) Create(_ context.Context, s *entity.Seat) error { f.created = s; return f.createErr }
func (f *fakeRepo) Update(_ context.Context, s *entity.Seat) error { f.updated = s; return f.updateErr }
func (f *fakeRepo) Delete(_ context.Context, id string) error      { f.deleted = id; return f.deleteErr }
func (f *fakeRepo) ExistsBySeatPosition(_ context.Context, _, num, row, excl string) (bool, error) {
	f.gotNum, f.gotRow, f.gotExcl = num, row, excl
	return f.exists, f.existsErr
}

func newBiz(r *fakeRepo) *business {
	return &business{repository: r, cache: fakeCache{}, roCache: fakeCache{}, redisClient: fakeRedis{}}
}

var ctx = context.Background()

func sp(s string) *string { return &s }

func base() *entity.Seat {
	return &entity.Seat{Id: "s1", RoomId: "r1", SeatNumber: "01", RowNumber: "A", SeatType: entity.SeatTypeRegular, Status: entity.SeatStatusAvailable}
}

func TestCreateSeat(t *testing.T) {
	t.Run("TC-SEAT-BIZ-01a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).CreateSeat(ctx, nil); !errors.Is(err, ErrInvalidSeatData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-01b", func(t *testing.T) { // empty seat number
		r := &fakeRepo{}
		err := newBiz(r).CreateSeat(ctx, &entity.Seat{RoomId: "r1", RowNumber: "A", SeatNumber: ""})
		if !errors.Is(err, ErrInvalidSeatData) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-01c", func(t *testing.T) { // duplicate position
		r := &fakeRepo{exists: true}
		err := newBiz(r).CreateSeat(ctx, base())
		if !errors.Is(err, ErrSeatPositionExists) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-01d", func(t *testing.T) {
		boom := errors.New("db")
		if err := newBiz(&fakeRepo{existsErr: boom}).CreateSeat(ctx, base()); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-01e", func(t *testing.T) {
		boom := errors.New("fk room")
		if err := newBiz(&fakeRepo{createErr: boom}).CreateSeat(ctx, base()); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-01f", func(t *testing.T) {
		r := &fakeRepo{}
		if err := newBiz(r).CreateSeat(ctx, base()); err != nil || r.created == nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-01g", func(t *testing.T) { // seat_type must be one of REGULAR/VIP/COUPLE (expected behaviour)
		r := &fakeRepo{}
		s := base()
		s.SeatType = entity.SeatType("BOGUS")
		if err := newBiz(r).CreateSeat(ctx, s); err == nil {
			t.Fatalf("seat_type BOGUS was accepted and persisted")
		}
	})
}

func TestUpdateSeat(t *testing.T) {
	t.Run("TC-SEAT-BIZ-02a", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).UpdateSeat(ctx, "s1", &entity.UpdateSeatRequest{})
		if !errors.Is(err, ErrSeatNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-02b", func(t *testing.T) { // seat number changed -> position check excludes itself
		r := &fakeRepo{seat: base()}
		if err := newBiz(r).UpdateSeat(ctx, "s1", &entity.UpdateSeatRequest{SeatNumber: sp("Z2")}); err != nil {
			t.Fatal(err)
		}
		if r.gotNum != "Z2" || r.gotRow != "A" || r.gotExcl != "s1" || r.updated.SeatNumber != "Z2" {
			t.Fatalf("num=%s row=%s excl=%s updated=%+v", r.gotNum, r.gotRow, r.gotExcl, r.updated)
		}
	})
	t.Run("TC-SEAT-BIZ-02c", func(t *testing.T) { // position taken
		r := &fakeRepo{seat: base(), exists: true}
		err := newBiz(r).UpdateSeat(ctx, "s1", &entity.UpdateSeatRequest{RowNumber: sp("B")})
		if !errors.Is(err, ErrSeatPositionExists) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-02d", func(t *testing.T) { // empty seat number rejected as invalid data
		r := &fakeRepo{seat: base()}
		err := newBiz(r).UpdateSeat(ctx, "s1", &entity.UpdateSeatRequest{SeatNumber: sp("")})
		if !errors.Is(err, ErrInvalidSeatData) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-02e", func(t *testing.T) { // type + status only: no position check
		r := &fakeRepo{seat: base()}
		ty, st := entity.SeatTypeVIP, entity.SeatStatusMaintenance
		if err := newBiz(r).UpdateSeat(ctx, "s1", &entity.UpdateSeatRequest{SeatType: &ty, Status: &st}); err != nil {
			t.Fatal(err)
		}
		if r.gotNum != "" || r.updated.SeatType != ty || r.updated.Status != st {
			t.Fatalf("updated=%+v gotNum=%q", r.updated, r.gotNum)
		}
	})
	t.Run("TC-SEAT-BIZ-02f", func(t *testing.T) {
		boom := errors.New("db")
		r := &fakeRepo{seat: base(), updateErr: boom}
		if err := newBiz(r).UpdateSeat(ctx, "s1", &entity.UpdateSeatRequest{}); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestDeleteSeat(t *testing.T) {
	t.Run("TC-SEAT-BIZ-03a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).DeleteSeat(ctx, "s1"); !errors.Is(err, ErrSeatNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-03b", func(t *testing.T) {
		r := &fakeRepo{seat: base()}
		if err := newBiz(r).DeleteSeat(ctx, "s1"); err != nil || r.deleted != "s1" {
			t.Fatalf("err=%v deleted=%q", err, r.deleted)
		}
	})
	t.Run("TC-SEAT-BIZ-03c", func(t *testing.T) {
		boom := errors.New("fk booking")
		if err := newBiz(&fakeRepo{seat: base(), deleteErr: boom}).DeleteSeat(ctx, "s1"); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestUpdateSeatStatus(t *testing.T) {
	t.Run("TC-SEAT-BIZ-04a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).UpdateSeatStatus(ctx, "s1", entity.SeatStatusMaintenance); !errors.Is(err, ErrSeatNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-04b", func(t *testing.T) { // AVAILABLE -> MAINTENANCE
		r := &fakeRepo{seat: base()}
		if err := newBiz(r).UpdateSeatStatus(ctx, "s1", entity.SeatStatusMaintenance); err != nil {
			t.Fatal(err)
		}
		if r.updated.Status != entity.SeatStatusMaintenance {
			t.Fatalf("status=%s", r.updated.Status)
		}
	})
	t.Run("TC-SEAT-BIZ-04c", func(t *testing.T) { // unknown status must be rejected (expected behaviour)
		r := &fakeRepo{seat: base()}
		if err := newBiz(r).UpdateSeatStatus(ctx, "s1", entity.SeatStatus("FOO")); err == nil {
			t.Fatalf("status FOO was accepted and persisted")
		}
	})
}

func TestGetSeats(t *testing.T) {
	t.Run("TC-SEAT-BIZ-05a", func(t *testing.T) {
		for _, pr := range [][2]int{{0, 10}, {1, 0}} {
			if _, _, err := newBiz(&fakeRepo{}).GetSeats(ctx, pr[0], pr[1], "", "", "", "", ""); !errors.Is(err, ErrInvalidSeatData) {
				t.Fatalf("page=%d size=%d err=%v", pr[0], pr[1], err)
			}
		}
	})
	t.Run("TC-SEAT-BIZ-05b", func(t *testing.T) {
		r := &fakeRepo{}
		seats, total, err := newBiz(r).GetSeats(ctx, 2, 50, "", "r1", "", "", "")
		if err != nil || len(seats) != 1 || total != 1 || r.gotOffset != 50 {
			t.Fatalf("seats=%v total=%d err=%v offset=%d", seats, total, err, r.gotOffset)
		}
	})
	t.Run("TC-SEAT-BIZ-05c", func(t *testing.T) {
		if _, err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).GetSeatById(ctx, "x"); !errors.Is(err, ErrSeatNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-SEAT-BIZ-05d", func(t *testing.T) {
		s, err := newBiz(&fakeRepo{seat: base()}).GetSeatById(ctx, "s1")
		if err != nil || s.Id != "s1" {
			t.Fatalf("s=%v err=%v", s, err)
		}
	})
}

