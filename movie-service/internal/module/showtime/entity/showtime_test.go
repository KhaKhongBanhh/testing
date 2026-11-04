package entity

import (
	"testing"
	"time"
)

var base = time.Date(2026, 5, 1, 14, 0, 0, 0, time.UTC)

func at(h, m, s, ns int) time.Time {
	return time.Date(2026, 5, 1, h, m, s, ns, time.UTC)
}

func TestTruncateToHalfHour(t *testing.T) {
	cases := []struct {
		id        string
		in, want  time.Time
	}{
		{"TC-SHOWTIME-ENTITY-01a", at(14, 20, 0, 0), at(14, 0, 0, 0)},
		{"TC-SHOWTIME-ENTITY-01b", at(14, 29, 0, 0), at(14, 0, 0, 0)},
		{"TC-SHOWTIME-ENTITY-01c", at(14, 30, 0, 0), at(14, 30, 0, 0)},
		{"TC-SHOWTIME-ENTITY-01d", at(14, 55, 0, 0), at(14, 30, 0, 0)},
		{"TC-SHOWTIME-ENTITY-01e", at(10, 15, 45, 123456789), at(10, 0, 0, 0)},
		{"TC-SHOWTIME-ENTITY-01f", at(14, 0, 0, 0), at(14, 0, 0, 0)},
		{"TC-SHOWTIME-ENTITY-01g", at(23, 59, 59, 0), at(23, 30, 0, 0)},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := TruncateToHalfHour(c.in); !got.Equal(c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func validShowtime() Showtime {
	return Showtime{MovieId: "m1", RoomId: "r1", BasePrice: 80000, StartTime: base, EndTime: base.Add(2 * time.Hour)}
}

func TestShowtimeIsValid(t *testing.T) {
	mod := func(f func(s *Showtime)) Showtime { s := validShowtime(); f(&s); return s }
	cases := []struct {
		id   string
		st   Showtime
		want bool
	}{
		{"TC-SHOWTIME-ENTITY-02a", validShowtime(), true},
		{"TC-SHOWTIME-ENTITY-02b", mod(func(s *Showtime) { s.MovieId = "" }), false},
		{"TC-SHOWTIME-ENTITY-02c", mod(func(s *Showtime) { s.RoomId = "" }), false},
		{"TC-SHOWTIME-ENTITY-02d", mod(func(s *Showtime) { s.BasePrice = -1 }), false},
		{"TC-SHOWTIME-ENTITY-02e", mod(func(s *Showtime) { s.EndTime = s.StartTime }), false},
		{"TC-SHOWTIME-ENTITY-02f", mod(func(s *Showtime) { s.EndTime = s.StartTime.Add(-time.Minute) }), false},
		{"TC-SHOWTIME-ENTITY-02g", mod(func(s *Showtime) { s.BasePrice = 0 }), true}, // entity allows 0; REST binding (required) rejects it
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.st.IsValid(); got != c.want {
				t.Fatalf("IsValid()=%v, want %v", got, c.want)
			}
		})
	}
}

func TestShowtimeIsActiveStatus(t *testing.T) {
	cases := []struct {
		id   string
		st   ShowtimeStatus
		want bool
	}{
		{"TC-SHOWTIME-ENTITY-03a", ShowtimeStatusScheduled, true},
		{"TC-SHOWTIME-ENTITY-03b", ShowtimeStatusOngoing, true},
		{"TC-SHOWTIME-ENTITY-03c", ShowtimeStatusCompleted, false},
		{"TC-SHOWTIME-ENTITY-03d", ShowtimeStatusCanceled, false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			s := &Showtime{Status: c.st}
			if got := s.IsActiveStatus(); got != c.want {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestShowtimeIsUpcoming(t *testing.T) {
	now := time.Now()
	cases := []struct {
		id   string
		st   Showtime
		want bool
	}{
		{"TC-SHOWTIME-ENTITY-04a", Showtime{Status: ShowtimeStatusScheduled, StartTime: now.Add(time.Hour)}, true},
		{"TC-SHOWTIME-ENTITY-04b", Showtime{Status: ShowtimeStatusScheduled, StartTime: now.Add(-time.Hour)}, false},
		{"TC-SHOWTIME-ENTITY-04c", Showtime{Status: ShowtimeStatusOngoing, StartTime: now.Add(time.Hour)}, false},
		{"TC-SHOWTIME-ENTITY-04d", Showtime{Status: ShowtimeStatusCanceled, StartTime: now.Add(time.Hour)}, false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.st.IsUpcoming(); got != c.want {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestShowtimeIsOngoing(t *testing.T) {
	now := time.Now()
	cases := []struct {
		id   string
		st   Showtime
		want bool
	}{
		{"TC-SHOWTIME-ENTITY-05a", Showtime{Status: ShowtimeStatusOngoing}, true},
		{"TC-SHOWTIME-ENTITY-05b", Showtime{Status: ShowtimeStatusScheduled, StartTime: now.Add(-time.Hour), EndTime: now.Add(time.Hour)}, true},
		{"TC-SHOWTIME-ENTITY-05c", Showtime{Status: ShowtimeStatusCompleted}, false},
		{"TC-SHOWTIME-ENTITY-05d", Showtime{Status: ShowtimeStatusScheduled, StartTime: now.Add(time.Hour), EndTime: now.Add(3 * time.Hour)}, false},
		{"TC-SHOWTIME-ENTITY-05e", Showtime{Status: ShowtimeStatusScheduled, StartTime: now.Add(-3 * time.Hour), EndTime: now.Add(-time.Hour)}, false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.st.IsOngoing(); got != c.want {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestShowtimeCalculateDuration(t *testing.T) {
	t.Run("TC-SHOWTIME-ENTITY-06a", func(t *testing.T) {
		s := &Showtime{StartTime: at(14, 0, 0, 0), EndTime: at(16, 30, 0, 0)}
		if got := s.CalculateDuration(); got != 150*time.Minute {
			t.Fatalf("got %v", got)
		}
	})
}

func validReq() CreateShowtimeRequest {
	return CreateShowtimeRequest{MovieId: "m1", RoomId: "r1", StartTime: base, EndTime: base.Add(2 * time.Hour), Format: ShowtimeFormat2D, BasePrice: 80000}
}

func TestCreateShowtimeRequestIsValid(t *testing.T) {
	mod := func(f func(r *CreateShowtimeRequest)) CreateShowtimeRequest { r := validReq(); f(&r); return r }
	cases := []struct {
		id   string
		req  CreateShowtimeRequest
		want bool
	}{
		{"TC-SHOWTIME-ENTITY-07a", validReq(), true},
		{"TC-SHOWTIME-ENTITY-07b", mod(func(r *CreateShowtimeRequest) { r.MovieId = "" }), false},
		{"TC-SHOWTIME-ENTITY-07c", mod(func(r *CreateShowtimeRequest) { r.BasePrice = -100 }), false},
		{"TC-SHOWTIME-ENTITY-07d", mod(func(r *CreateShowtimeRequest) { r.EndTime = r.StartTime }), false},
		{"TC-SHOWTIME-ENTITY-07e", mod(func(r *CreateShowtimeRequest) { r.RoomId = "" }), false},
		{"TC-SHOWTIME-ENTITY-07f", mod(func(r *CreateShowtimeRequest) { r.EndTime = r.StartTime.Add(-time.Hour) }), false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.req.IsValid(); got != c.want {
				t.Fatalf("IsValid()=%v, want %v", got, c.want)
			}
		})
	}
}

func TestCreateShowtimeRequestToShowtime(t *testing.T) {
	t.Run("TC-SHOWTIME-ENTITY-08a", func(t *testing.T) {
		r := validReq()
		r.StartTime, r.EndTime = at(14, 20, 0, 0), at(16, 45, 0, 0)
		s := r.ToShowtime()
		if !s.StartTime.Equal(at(14, 0, 0, 0)) || !s.EndTime.Equal(at(16, 30, 0, 0)) || s.Status != ShowtimeStatusScheduled {
			t.Fatalf("got %v - %v (%s)", s.StartTime, s.EndTime, s.Status)
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-08b", func(t *testing.T) {
		r := validReq()
		r.StartTime, r.EndTime = at(14, 20, 0, 0), at(14, 10, 0, 0)
		s := r.ToShowtime()
		if !s.StartTime.Equal(at(14, 0, 0, 0)) || !s.EndTime.Equal(at(16, 0, 0, 0)) {
			t.Fatalf("got %v - %v, want 14:00 - 16:00", s.StartTime, s.EndTime)
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-08c", func(t *testing.T) {
		r := validReq()
		r.StartTime, r.EndTime = at(14, 0, 0, 0), at(12, 0, 0, 0) // end before start
		s := r.ToShowtime()
		if !s.EndTime.Equal(at(16, 0, 0, 0)) {
			t.Fatalf("end=%v, want 16:00", s.EndTime)
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-08d", func(t *testing.T) {
		r := validReq()
		s := r.ToShowtime()
		if s.MovieId != "m1" || s.RoomId != "r1" || s.Format != ShowtimeFormat2D || s.BasePrice != 80000 {
			t.Fatalf("fields not copied: %+v", s)
		}
	})
}

func TestShowtimeResponses(t *testing.T) {
	t.Run("TC-SHOWTIME-ENTITY-09a", func(t *testing.T) {
		s := &Showtime{Id: "st1", MovieId: "m1", RoomId: "r1", StartTime: at(14, 0, 0, 0), EndTime: at(16, 0, 0, 0),
			Format: ShowtimeFormat3D, BasePrice: 100, Status: ShowtimeStatusScheduled, CreatedAt: base}
		r := ToShowtimeResponse(s)
		if r.Id != "st1" || r.StartTime != "2026-05-01T14:00:00Z" || r.Duration != "2h0m0s" || r.UpdatedAt != nil || r.Format != ShowtimeFormat3D {
			t.Fatalf("bad mapping: %+v", r)
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-09b", func(t *testing.T) {
		if ToShowtimeResponse(nil) != nil {
			t.Fatalf("want nil")
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-09c", func(t *testing.T) {
		ts := base
		r := ToShowtimeResponse(&Showtime{Id: "st1", StartTime: base, EndTime: base.Add(time.Hour), UpdatedAt: &ts})
		if r.UpdatedAt == nil || *r.UpdatedAt != "2026-05-01T14:00:00Z" {
			t.Fatalf("UpdatedAt=%v", r.UpdatedAt)
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-09d", func(t *testing.T) {
		if ToShowtimeBookingResponse(nil) != nil {
			t.Fatalf("want nil")
		}
		r := ToShowtimeBookingResponse(&Showtime{Id: "st1", StartTime: base, EndTime: base.Add(time.Hour), Seats: []*Seat{{Id: "s1"}}})
		if r.Id != "st1" || len(r.Seats) != 1 {
			t.Fatalf("bad mapping: %+v", r)
		}
	})
	t.Run("TC-SHOWTIME-ENTITY-10a", func(t *testing.T) {
		r := ToShowtimesResponse([]*Showtime{{Id: "a", StartTime: base, EndTime: base.Add(time.Hour)}}, 1, 10, 11)
		if len(r.Data) != 1 || r.Paging.TotalPages != 2 {
			t.Fatalf("resp=%+v paging=%+v", r, r.Paging)
		}
	})
}

