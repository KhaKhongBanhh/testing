package entity

import "testing"

func TestSeatIsValid(t *testing.T) {
	cases := []struct {
		id   string
		seat Seat
		want bool
	}{
		{"TC-SEAT-ENTITY-01a", Seat{RoomId: "r1", SeatNumber: "01", RowNumber: "A"}, true},
		{"TC-SEAT-ENTITY-01b", Seat{RoomId: "", SeatNumber: "01", RowNumber: "A"}, false},
		{"TC-SEAT-ENTITY-01c", Seat{RoomId: "r1", SeatNumber: "", RowNumber: "A"}, false},
		{"TC-SEAT-ENTITY-01d", Seat{RoomId: "r1", SeatNumber: "01", RowNumber: ""}, false},
		{"TC-SEAT-ENTITY-01e", Seat{}, false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.seat.IsValid(); got != c.want {
				t.Fatalf("IsValid()=%v, want %v", got, c.want)
			}
		})
	}
}

func TestGetSeatTypePriceMultiplier(t *testing.T) {
	cases := []struct {
		id   string
		typ  SeatType
		want float64
	}{
		{"TC-SEAT-ENTITY-02a", SeatTypeRegular, 1},
		{"TC-SEAT-ENTITY-02b", SeatTypeVIP, 1.5},
		{"TC-SEAT-ENTITY-02c", SeatTypeCouple, 2.5},
		{"TC-SEAT-ENTITY-02d", SeatType("UNKNOWN"), 1},
		{"TC-SEAT-ENTITY-02e", SeatType(""), 1},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := GetSeatTypePriceMultiplier(c.typ); got != c.want {
				t.Fatalf("multiplier=%v, want %v", got, c.want)
			}
		})
	}
}

func TestSeatCalculatePrice(t *testing.T) {
	cases := []struct {
		id    string
		typ   SeatType
		base  float64
		want  float64
	}{
		{"TC-SEAT-ENTITY-03a", SeatTypeRegular, 100000, 100000},
		{"TC-SEAT-ENTITY-03b", SeatTypeVIP, 100000, 150000},
		{"TC-SEAT-ENTITY-03c", SeatTypeCouple, 100000, 250000},
		{"TC-SEAT-ENTITY-03d", SeatTypeVIP, 0, 0},
		{"TC-SEAT-ENTITY-03e", SeatType("UNKNOWN"), 80000, 80000},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			s := &Seat{SeatType: c.typ}
			if got := s.CalculatePrice(c.base); got != c.want {
				t.Fatalf("price=%v, want %v", got, c.want)
			}
		})
	}
}

