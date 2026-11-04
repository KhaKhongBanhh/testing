package entity

import "testing"

func TestCreateMovieRequestToEntity(t *testing.T) {
	t.Run("TC-MOVIE-ENTITY-03a", func(t *testing.T) {
		m := (&CreateMovieRequest{Title: "A", Duration: 100}).ToEntity()
		if m.Status != MovieStatusUpcoming {
			t.Fatalf("status=%s, want UPCOMING", m.Status)
		}
	})
	t.Run("TC-MOVIE-ENTITY-03b", func(t *testing.T) {
		m := (&CreateMovieRequest{Title: "A", Duration: 100, Status: "SHOWING"}).ToEntity()
		if m.Status != MovieStatusShowing {
			t.Fatalf("status=%s, want SHOWING", m.Status)
		}
	})
	t.Run("TC-MOVIE-ENTITY-03c", func(t *testing.T) {
		m := (&CreateMovieRequest{Title: "T", Director: "D", Cast: "C", Duration: 90, Description: "desc", TrailerURL: "t", PosterURL: "p"}).ToEntity()
		if m.Title != "T" || m.Director != "D" || m.Cast != "C" || m.Duration != 90 || m.Description != "desc" || m.TrailerURL != "t" || m.PosterURL != "p" {
			t.Fatalf("fields not copied: %+v", m)
		}
	})
}

func TestUpdateMovieRequestToEntity(t *testing.T) {
	t.Run("TC-MOVIE-ENTITY-03d", func(t *testing.T) {
		m := (&UpdateMovieRequest{Title: "A", Duration: 100}).ToEntity("id-1")
		if m.Id != "id-1" {
			t.Fatalf("id=%s", m.Id)
		}
		// No default is applied on update: an omitted status stays empty.
		if m.Status != "" {
			t.Fatalf("status=%q, want empty", m.Status)
		}
	})
	t.Run("TC-MOVIE-ENTITY-03e", func(t *testing.T) {
		m := (&UpdateMovieRequest{Title: "A", Duration: 100, Status: "ENDED"}).ToEntity("id-1")
		if m.Status != MovieStatusEnded {
			t.Fatalf("status=%s, want ENDED", m.Status)
		}
	})
}

func TestToMovieResponse(t *testing.T) {
	t.Run("TC-MOVIE-ENTITY-04a", func(t *testing.T) {
		resp := ToMovieResponse(&Movie{Id: "m-1", Title: "Test", Status: MovieStatusShowing, Duration: 90})
		if resp.Id != "m-1" || resp.Title != "Test" || resp.Status != "SHOWING" || resp.Duration != 90 {
			t.Fatalf("bad mapping: %+v", resp)
		}
	})
	t.Run("TC-MOVIE-ENTITY-04b", func(t *testing.T) {
		resp := ToMovieResponse(&Movie{Id: "m-1", MovieGenres: []*MovieGenre{
			{Genre: &Genre{Name: "Action"}}, {Genre: &Genre{Name: "Drama"}},
		}})
		if len(resp.Genres) != 2 || resp.Genres[0] != "Action" || resp.Genres[1] != "Drama" {
			t.Fatalf("genres=%v", resp.Genres)
		}
	})
	t.Run("TC-MOVIE-ENTITY-04c", func(t *testing.T) {
		resp := ToMovieResponse(&Movie{Id: "m-1"})
		if resp.Genres == nil || len(resp.Genres) != 0 {
			t.Fatalf("genres=%v, want empty non-nil slice", resp.Genres)
		}
	})
}

func TestToMoviesResponse(t *testing.T) {
	cases := []struct {
		id          string
		total, size int
		wantPages   int
	}{
		{"TC-MOVIE-ENTITY-05a", 25, 10, 3},
		{"TC-MOVIE-ENTITY-05b", 0, 10, 1},
		{"TC-MOVIE-ENTITY-05c", 30, 10, 3},  // exact multiple
		{"TC-MOVIE-ENTITY-05d", 1, 10, 1},   // boundary
		{"TC-MOVIE-ENTITY-05e", 11, 10, 2},  // boundary: one over a full page
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			resp := ToMoviesResponse(nil, 1, c.size, c.total)
			if resp.Meta.TotalPages != c.wantPages {
				t.Fatalf("TotalPages=%d, want %d", resp.Meta.TotalPages, c.wantPages)
			}
			if resp.Meta.Total != c.total || resp.Meta.Size != c.size || resp.Meta.Page != 1 {
				t.Fatalf("meta=%+v", resp.Meta)
			}
		})
	}
	t.Run("TC-MOVIE-ENTITY-05f", func(t *testing.T) {
		resp := ToMoviesResponse([]*Movie{{Id: "a"}, {Id: "b"}}, 1, 10, 2)
		if len(resp.Movies) != 2 || resp.Movies[1].Id != "b" {
			t.Fatalf("movies=%v", resp.Movies)
		}
	})
}

