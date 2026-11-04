package entity

import "testing"

func TestMovieIsValid(t *testing.T) {
	cases := []struct {
		id    string
		movie Movie
		want  bool
	}{
		{"TC-MOVIE-ENTITY-01a", Movie{Title: "Avengers", Duration: 120}, true},
		{"TC-MOVIE-ENTITY-01b", Movie{Title: "", Duration: 120}, false},
		{"TC-MOVIE-ENTITY-01c", Movie{Title: "Test", Duration: 0}, false},
		{"TC-MOVIE-ENTITY-01d", Movie{Title: "Test", Duration: -1}, false},
		{"TC-MOVIE-ENTITY-01e", Movie{Title: "Test", Duration: 1}, true}, // boundary: smallest valid duration
		{"TC-MOVIE-ENTITY-01f", Movie{Title: "", Duration: 0}, false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			if got := c.movie.IsValid(); got != c.want {
				t.Fatalf("IsValid()=%v, want %v", got, c.want)
			}
		})
	}
}

func TestMovieCanTransitionTo(t *testing.T) {
	cases := []struct {
		id        string
		cur, next MovieStatus
		want      bool
	}{
		{"TC-MOVIE-ENTITY-02a", MovieStatusUpcoming, MovieStatusShowing, true},
		{"TC-MOVIE-ENTITY-02b", MovieStatusUpcoming, MovieStatusEnded, false},
		{"TC-MOVIE-ENTITY-02c", MovieStatusShowing, MovieStatusEnded, true},
		{"TC-MOVIE-ENTITY-02d", MovieStatusShowing, MovieStatusUpcoming, false},
		{"TC-MOVIE-ENTITY-02e", MovieStatusEnded, MovieStatusUpcoming, false},
		{"TC-MOVIE-ENTITY-02f", MovieStatusEnded, MovieStatusShowing, false},
		{"TC-MOVIE-ENTITY-02g", MovieStatusUpcoming, MovieStatusUpcoming, false},
		{"TC-MOVIE-ENTITY-02h", MovieStatusShowing, MovieStatusShowing, false},
		{"TC-MOVIE-ENTITY-02i", MovieStatusEnded, MovieStatusEnded, false},
		{"TC-MOVIE-ENTITY-02j", MovieStatusUpcoming, MovieStatus(""), false},
		{"TC-MOVIE-ENTITY-02k", MovieStatusUpcoming, MovieStatus("showing"), false}, // lowercase is not normalised
		{"TC-MOVIE-ENTITY-02l", MovieStatus("BOGUS"), MovieStatusShowing, false},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			m := &Movie{Status: c.cur}
			if got := m.CanTransitionTo(c.next); got != c.want {
				t.Fatalf("%s -> %s = %v, want %v", c.cur, c.next, got, c.want)
			}
		})
	}
}

