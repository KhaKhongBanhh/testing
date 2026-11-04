package entity

import "testing"

func TestCreateSeatRequestToSeat(t *testing.T) {
	t.Run("TC-SEAT-ENTITY-04a", func(t *testing.T) {
		seat := (&CreateSeatRequest{RoomId: "r1", SeatNumber: "01", RowNumber: "A", SeatType: SeatTypeVIP}).ToSeat()
		if seat.Status != SeatStatusAvailable {
			t.Fatalf("status=%s, want AVAILABLE", seat.Status)
		}
	})
	t.Run("TC-SEAT-ENTITY-04b", func(t *testing.T) {
		seat := (&CreateSeatRequest{RoomId: "r1", SeatNumber: "01", RowNumber: "A", SeatType: SeatTypeVIP}).ToSeat()
		if seat.RoomId != "r1" || seat.SeatNumber != "01" || seat.RowNumber != "A" || seat.SeatType != SeatTypeVIP {
			t.Fatalf("fields not copied: %+v", seat)
		}
	})
}

func TestToSeatResponse(t *testing.T) {
	t.Run("TC-SEAT-ENTITY-05a", func(t *testing.T) {
		resp := ToSeatResponse(&Seat{Id: "s-1", SeatType: SeatTypeCouple, Status: SeatStatusMaintenance})
		if resp.Id != "s-1" || resp.SeatType != SeatTypeCouple || resp.Status != SeatStatusMaintenance {
			t.Fatalf("bad mapping: %+v", resp)
		}
	})
	t.Run("TC-SEAT-ENTITY-05b", func(t *testing.T) {
		if ToSeatResponse(nil) != nil {
			t.Fatalf("want nil")
		}
	})
}

func TestToSeatsResponse(t *testing.T) {
	t.Run("TC-SEAT-ENTITY-06a", func(t *testing.T) {
		resp := ToSeatsResponse([]*Seat{{Id: "a"}, {Id: "b"}}, 1, 10, 2)
		if len(resp.Data) != 2 || resp.Paging == nil {
			t.Fatalf("resp=%+v", resp)
		}
	})
	t.Run("TC-SEAT-ENTITY-06b", func(t *testing.T) {
		resp := ToSeatsResponse(nil, 2, 50, 120)
		if resp.Paging.TotalPages != 3 || resp.Paging.Page != 2 {
			t.Fatalf("paging=%+v", resp.Paging)
		}
	})
}

