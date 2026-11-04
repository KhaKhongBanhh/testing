package business

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"movie-service/internal/module/room/entity"

	rediscache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type fakeCache struct{}

func (fakeCache) Get(context.Context, string, any) error { return rediscache.ErrCacheMiss }
func (fakeCache) Set(context.Context, string, any, time.Duration) error { return nil }
func (fakeCache) Delete(context.Context, string) error                  { return nil }

type fakeRepo struct {
	room      *entity.Room
	getErr    error
	exists    bool
	existsErr error
	createErr error
	updateErr error
	deleteErr error
	created   *entity.Room
	updated   *entity.Room
	deleted   string
	gotExcl   string
	gotLimit  int
	gotOffset int
}

func (f *fakeRepo) GetByID(context.Context, string) (*entity.Room, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	cp := *f.room
	return &cp, nil
}
func (f *fakeRepo) GetMany(_ context.Context, limit, offset int, _ string, _ entity.RoomType, _ entity.RoomStatus) ([]*entity.Room, error) {
	f.gotLimit, f.gotOffset = limit, offset
	return []*entity.Room{{Id: "a"}}, nil
}
func (f *fakeRepo) GetTotalCount(context.Context, string, entity.RoomType, entity.RoomStatus) (int, error) {
	return 1, nil
}
func (f *fakeRepo) Create(_ context.Context, r *entity.Room) error { f.created = r; return f.createErr }
func (f *fakeRepo) Update(_ context.Context, r *entity.Room) error { f.updated = r; return f.updateErr }
func (f *fakeRepo) Delete(_ context.Context, id string) error      { f.deleted = id; return f.deleteErr }
func (f *fakeRepo) ExistsByRoomNumber(_ context.Context, _ int, excl string) (bool, error) {
	f.gotExcl = excl
	return f.exists, f.existsErr
}

func newBiz(r *fakeRepo) *business {
	return &business{repository: r, cache: fakeCache{}, roCache: fakeCache{}, redisClient: nil}
}

var ctx = context.Background()

func intp(i int) *int { return &i }

func TestCreateRoom(t *testing.T) {
	t.Run("TC-ROOM-BIZ-01a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).CreateRoom(ctx, nil); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-01b", func(t *testing.T) { // capacity 0
		r := &fakeRepo{}
		err := newBiz(r).CreateRoom(ctx, &entity.Room{RoomNumber: 1, Capacity: 0})
		if !errors.Is(err, ErrInvalidRoomData) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-01c", func(t *testing.T) { // capacity negative
		if err := newBiz(&fakeRepo{}).CreateRoom(ctx, &entity.Room{RoomNumber: 1, Capacity: -1}); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-01d", func(t *testing.T) { // duplicate room number
		r := &fakeRepo{exists: true}
		err := newBiz(r).CreateRoom(ctx, &entity.Room{RoomNumber: 999, Capacity: 100})
		if !errors.Is(err, ErrRoomNumberExists) || r.created != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-01e", func(t *testing.T) { // exists check error
		boom := errors.New("db")
		err := newBiz(&fakeRepo{existsErr: boom}).CreateRoom(ctx, &entity.Room{RoomNumber: 1, Capacity: 1})
		if !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-01f", func(t *testing.T) { // create error
		boom := errors.New("db")
		err := newBiz(&fakeRepo{createErr: boom}).CreateRoom(ctx, &entity.Room{RoomNumber: 1, Capacity: 1})
		if !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-01g", func(t *testing.T) { // success
		r := &fakeRepo{}
		if err := newBiz(r).CreateRoom(ctx, &entity.Room{RoomNumber: 999, Capacity: 100}); err != nil || r.created == nil {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestUpdateRoom(t *testing.T) {
	base := func() *entity.Room {
		return &entity.Room{Id: "r1", RoomNumber: 5, Capacity: 100, RoomType: entity.RoomTypeStandard, Status: entity.RoomStatusActive}
	}
	t.Run("TC-ROOM-BIZ-02a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).UpdateRoom(ctx, "", &entity.UpdateRoomRequest{}); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-02b", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).UpdateRoom(ctx, "r1", nil); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-02c", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).UpdateRoom(ctx, "r1", &entity.UpdateRoomRequest{})
		if !errors.Is(err, ErrRoomNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-02d", func(t *testing.T) { // room number already used by another room
		r := &fakeRepo{room: base(), exists: true}
		err := newBiz(r).UpdateRoom(ctx, "r1", &entity.UpdateRoomRequest{RoomNumber: intp(1)})
		if !errors.Is(err, ErrRoomNumberExists) || r.gotExcl != "r1" {
			t.Fatalf("err=%v excl=%q", err, r.gotExcl)
		}
	})
	t.Run("TC-ROOM-BIZ-02e", func(t *testing.T) { // capacity update
		r := &fakeRepo{room: base()}
		if err := newBiz(r).UpdateRoom(ctx, "r1", &entity.UpdateRoomRequest{Capacity: intp(150)}); err != nil {
			t.Fatal(err)
		}
		if r.updated.Capacity != 150 || r.updated.RoomNumber != 5 {
			t.Fatalf("updated=%+v", r.updated)
		}
	})
	t.Run("TC-ROOM-BIZ-02f", func(t *testing.T) { // capacity 0 rejected
		r := &fakeRepo{room: base()}
		err := newBiz(r).UpdateRoom(ctx, "r1", &entity.UpdateRoomRequest{Capacity: intp(0)})
		if !errors.Is(err, ErrInvalidRoomData) || r.updated != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-02g", func(t *testing.T) { // status + type update
		r := &fakeRepo{room: base()}
		st, ty := entity.RoomStatusMaintenance, entity.RoomTypeIMAX
		if err := newBiz(r).UpdateRoom(ctx, "r1", &entity.UpdateRoomRequest{Status: &st, RoomType: &ty}); err != nil {
			t.Fatal(err)
		}
		if r.updated.Status != st || r.updated.RoomType != ty {
			t.Fatalf("updated=%+v", r.updated)
		}
	})
	t.Run("TC-ROOM-BIZ-02h", func(t *testing.T) { // repository error
		boom := errors.New("db")
		r := &fakeRepo{room: base(), updateErr: boom}
		if err := newBiz(r).UpdateRoom(ctx, "r1", &entity.UpdateRoomRequest{Capacity: intp(10)}); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestDeleteRoom(t *testing.T) {
	t.Run("TC-ROOM-BIZ-03a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).DeleteRoom(ctx, ""); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-03b", func(t *testing.T) {
		r := &fakeRepo{}
		if err := newBiz(r).DeleteRoom(ctx, "r1"); err != nil || r.deleted != "r1" {
			t.Fatalf("err=%v deleted=%q", err, r.deleted)
		}
	})
	t.Run("TC-ROOM-BIZ-03c", func(t *testing.T) {
		boom := errors.New("fk")
		if err := newBiz(&fakeRepo{deleteErr: boom}).DeleteRoom(ctx, "r1"); !errors.Is(err, boom) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestUpdateRoomStatus(t *testing.T) {
	t.Run("TC-ROOM-BIZ-04a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).UpdateRoomStatus(ctx, "", entity.RoomStatusInactive); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-04b", func(t *testing.T) {
		err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).UpdateRoomStatus(ctx, "r1", entity.RoomStatusInactive)
		if !errors.Is(err, ErrRoomNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-04c", func(t *testing.T) { // ACTIVE -> INACTIVE
		r := &fakeRepo{room: &entity.Room{Id: "r1", RoomNumber: 1, Capacity: 1, Status: entity.RoomStatusActive}}
		if err := newBiz(r).UpdateRoomStatus(ctx, "r1", entity.RoomStatusInactive); err != nil {
			t.Fatal(err)
		}
		if r.updated.Status != entity.RoomStatusInactive {
			t.Fatalf("status=%s", r.updated.Status)
		}
	})
	t.Run("TC-ROOM-BIZ-04d", func(t *testing.T) { // an unknown status value must be rejected (expected behaviour)
		r := &fakeRepo{room: &entity.Room{Id: "r1", RoomNumber: 1, Capacity: 1, Status: entity.RoomStatusActive}}
		err := newBiz(r).UpdateRoomStatus(ctx, "r1", entity.RoomStatus("FOO"))
		if err == nil {
			t.Fatalf("status FOO was accepted and persisted: %+v", r.updated)
		}
	})
}

func TestValidateRoomForShowtime(t *testing.T) {
	t.Run("TC-ROOM-BIZ-05a", func(t *testing.T) {
		if err := newBiz(&fakeRepo{}).ValidateRoomForShowtime(ctx, ""); !errors.Is(err, ErrInvalidRoomData) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-05b", func(t *testing.T) {
		if err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).ValidateRoomForShowtime(ctx, "r1"); !errors.Is(err, ErrRoomNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-05c", func(t *testing.T) {
		for _, s := range []entity.RoomStatus{entity.RoomStatusInactive, entity.RoomStatusMaintenance} {
			r := &fakeRepo{room: &entity.Room{Id: "r1", Status: s}}
			if err := newBiz(r).ValidateRoomForShowtime(ctx, "r1"); !errors.Is(err, ErrRoomNotActive) {
				t.Fatalf("status %s: err=%v", s, err)
			}
		}
	})
	t.Run("TC-ROOM-BIZ-05d", func(t *testing.T) {
		r := &fakeRepo{room: &entity.Room{Id: "r1", Status: entity.RoomStatusActive}}
		if err := newBiz(r).ValidateRoomForShowtime(ctx, "r1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestGetRooms(t *testing.T) {
	t.Run("TC-ROOM-BIZ-06a", func(t *testing.T) {
		for _, pr := range [][2]int{{0, 10}, {1, 0}, {-1, -1}} {
			if _, _, err := newBiz(&fakeRepo{}).GetRooms(ctx, pr[0], pr[1], "", "", ""); !errors.Is(err, ErrInvalidRoomData) {
				t.Fatalf("page=%d size=%d err=%v", pr[0], pr[1], err)
			}
		}
	})
	t.Run("TC-ROOM-BIZ-06b", func(t *testing.T) {
		r := &fakeRepo{}
		rooms, total, err := newBiz(r).GetRooms(ctx, 3, 10, "", "", "")
		if err != nil || len(rooms) != 1 || total != 1 || r.gotOffset != 20 || r.gotLimit != 10 {
			t.Fatalf("rooms=%v total=%d err=%v offset=%d", rooms, total, err, r.gotOffset)
		}
	})
}

func TestGetRoomById(t *testing.T) {
	t.Run("TC-ROOM-BIZ-07a", func(t *testing.T) {
		if _, err := newBiz(&fakeRepo{getErr: sql.ErrNoRows}).GetRoomById(ctx, "x"); !errors.Is(err, ErrRoomNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TC-ROOM-BIZ-07b", func(t *testing.T) {
		got, err := newBiz(&fakeRepo{room: &entity.Room{Id: "r1"}}).GetRoomById(ctx, "r1")
		if err != nil || got.Id != "r1" {
			t.Fatalf("got=%v err=%v", got, err)
		}
	})
}

var _ = redis.Nil

