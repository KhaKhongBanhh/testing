package paging

import "testing"

func TestNewPageInfo(t *testing.T) {
	cases := []struct {
		id                string
		page, size, total int
		wantPages         int
	}{
		{"TC-PAGING-01a", 1, 10, 25, 3},
		{"TC-PAGING-01b", 1, 10, 20, 2},
		{"TC-PAGING-01c", 1, 10, 0, 0},
		{"TC-PAGING-01d", 1, 10, 1, 1},
		{"TC-PAGING-01e", 1, 10, 11, 2},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			p := NewPageInfo(c.page, c.size, c.total)
			if p.TotalPages != c.wantPages || p.Page != c.page || p.Size != c.size || p.Total != c.total {
				t.Fatalf("got %+v, want TotalPages=%d", p, c.wantPages)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	t.Run("TC-PAGING-02a", func(t *testing.T) {
		if stringToInt("7") != 7 {
			t.Fatal("want 7")
		}
	})
	t.Run("TC-PAGING-02b", func(t *testing.T) {
		if stringToInt("abc") != 0 {
			t.Fatal("want 0 on parse error")
		}
	})
}

