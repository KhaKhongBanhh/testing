package entity

import (
	"testing"
	"time"
)

func TestRoomIsValid(t *testing.T) {
	cases := []struct {
		id   string
		room Room
		want bool
	}{
		{"TC-ROOM-ENTITY-01a", Room{RoomNumber: 1, Capacity: 100}, true},
		{"TC-ROOM-ENTITY-01b", Room{RoomNumber: 0, Capacity: 100}, false},
		{"TC-ROOM-ENTITY-01c", Room{RoomNumber: -1, Capacity: 100}, false},
		{"TC-ROOM-ENTITY-01d", Room{RoomNumber: 5, Capacity: 0}, false},
		{"TC-ROOM-ENTITY-01e", Room{RoomNumber: 5, Capacity: -50}, false},
		{"TC-ROOM-ENTITY-01f", Room{RoomNumber: 0, Capacity: 0}, false},
		{"TC-ROOM-ENTITY-01g", Room{RoomNumber: 1, Capacity: 1}, true}, // boundary: smallest valid values
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.room.IsValid(); got != c.want {
				t.Fatalf("IsValid()=%v, want %v", got, c.want)
			}
		})
	}
}

func TestCreateRoomRequestToRoom(t *testing.T) {
	t.Run("TC-ROOM-ENTITY-02a", func(t *testing.T) {
		room := (&CreateRoomRequest{RoomNumber: 3, Capacity: 80, RoomType: RoomTypeStandard}).ToRoom()
		if room.Status != RoomStatusActive {
			t.Fatalf("status=%s, want ACTIVE", room.Status)
		}
	})
	t.Run("TC-ROOM-ENTITY-02b", func(t *testing.T) {
		room := (&CreateRoomRequest{RoomNumber: 3, Capacity: 80, RoomType: RoomTypeIMAX}).ToRoom()
		if room.RoomNumber != 3 || room.Capacity != 80 || room.RoomType != RoomTypeIMAX {
			t.Fatalf("fields not copied: %+v", room)
		}
	})
}

func TestToRoomResponse(t *testing.T) {
	t.Run("TC-ROOM-ENTITY-03a", func(t *testing.T) {
		resp := ToRoomResponse(&Room{Id: "r-1", RoomNumber: 2, Capacity: 10, RoomType: RoomTypeVIP, Status: RoomStatusActive})
		if resp.Id != "r-1" || resp.RoomNumber != 2 || resp.Capacity != 10 || resp.RoomType != RoomTypeVIP || resp.Status != RoomStatusActive {
			t.Fatalf("bad mapping: %+v", resp)
		}
		if resp.UpdatedAt != nil {
			t.Fatalf("UpdatedAt should be nil")
		}
	})
	t.Run("TC-ROOM-ENTITY-03b", func(t *testing.T) {
		if ToRoomResponse(nil) != nil {
			t.Fatalf("want nil")
		}
	})
	t.Run("TC-ROOM-ENTITY-03c", func(t *testing.T) {
		ts := time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC)
		resp := ToRoomResponse(&Room{Id: "r-1", CreatedAt: ts, UpdatedAt: &ts})
		if resp.CreatedAt != "2026-05-01T10:30:00Z" || resp.UpdatedAt == nil || *resp.UpdatedAt != "2026-05-01T10:30:00Z" {
			t.Fatalf("bad time format: %v %v", resp.CreatedAt, resp.UpdatedAt)
		}
	})
}

func TestToRoomsResponse(t *testing.T) {
	t.Run("TC-ROOM-ENTITY-04a", func(t *testing.T) {
		resp := ToRoomsResponse([]*Room{{Id: "a"}, {Id: "b"}}, 1, 10, 2)
		if len(resp.Data) != 2 || resp.Paging == nil {
			t.Fatalf("resp=%+v", resp)
		}
	})
	t.Run("TC-ROOM-ENTITY-04b", func(t *testing.T) {
		resp := ToRoomsResponse(nil, 1, 10, 25)
		if resp.Paging.TotalPages != 3 {
			t.Fatalf("TotalPages=%d, want 3", resp.Paging.TotalPages)
		}
	})
	t.Run("TC-ROOM-ENTITY-04c", func(t *testing.T) {
		resp := ToRoomsResponse(nil, 1, 10, 0)
		if resp.Data == nil || len(resp.Data) != 0 {
			t.Fatalf("Data should be an empty non-nil slice")
		}
	})
}

