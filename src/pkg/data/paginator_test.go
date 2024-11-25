package data

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	t.Parallel()
	p := NewPaginator(2, 10)
	if p.Number != 2 || p.PerPage != 10 {
		t.Errorf("NewPaginator(2, 10) = %v; want {Number:2, PerPage:10}", p)
	}
}

func TestSetCount(t *testing.T) {
	tests := []struct {
		name        string
		paginator   *Paginator
		count       int
		wantPages   int
		wantNumber  int
		wantErr     bool
		expectedErr error
	}{
		{"Normal case", NewPaginator(1, 10), 100, 10, 1, false, nil},
		{"Zero count", NewPaginator(1, 10), 0, 1, 1, false, nil},
		{"Number too high", NewPaginator(11, 10), 100, 10, 11, true, ErrInvalidPage},
		{"Number too low", NewPaginator(-1, 10), 100, 10, 0, true, ErrInvalidPage},
		{"Number zero", NewPaginator(0, 10), 100, 10, 1, false, nil},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.paginator.SetCount(tt.count)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetCount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("SetCount() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			} else {
				if tt.paginator.NumPages != tt.wantPages {
					t.Errorf("NumPages = %v; want %v", tt.paginator.NumPages, tt.wantPages)
				}
				if tt.paginator.Number != tt.wantNumber {
					t.Errorf("Number = %v; want %v", tt.paginator.Number, tt.wantNumber)
				}
			}
		})
	}
}

func TestGetLimitOffset(t *testing.T) {
	tests := []struct {
		name       string
		paginator  *Paginator
		wantLimit  int
		wantOffset int
	}{
		{"First Page", NewPaginator(1, 10), 10, 0},
		{"Fifth Page", NewPaginator(5, 10), 10, 40},
		{"Tenth Page", NewPaginator(10, 10), 10, 90},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.paginator.SetCount(99)
			if err != nil {
				t.Fatalf("SetCount() error = %v", err)
			}
			gotLimit, gotOffset := tt.paginator.getLimitOffset()
			if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
				t.Errorf("GetLimitOffset() = (%v, %v); want (%v, %v)", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestHasNextPrevious(t *testing.T) {
	tests := []struct {
		name        string
		paginator   *Paginator
		wantHasNext bool
		wantHasPrev bool
	}{
		{"First Page", NewPaginator(1, 10), true, false},
		{"Middle Page", NewPaginator(5, 10), true, true},
		{"Last Page", NewPaginator(10, 10), false, true},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.paginator.SetCount(100)
			if err != nil {
				t.Fatalf("SetCount() error = %v", err)
			}
			if got := tt.paginator.HasNext(); got != tt.wantHasNext {
				t.Errorf("HasNext() = %v; want %v", got, tt.wantHasNext)
			}
			if got := tt.paginator.HasPrevious(); got != tt.wantHasPrev {
				t.Errorf("HasPrevious() = %v; want %v", got, tt.wantHasPrev)
			}
		})
	}
}

func TestPreviousCurrentNextStrings(t *testing.T) {
	t.Parallel()
	pager := NewPaginator(5, 10)
	err := pager.SetCount(100)
	if err != nil {
		t.Fatalf("SetCount() error = %v", err)
	}

	if got := pager.Previous().String(); got != "4" {
		t.Errorf("Previous() = %v; want 4", got)
	}
	if !pager.IsCurrent(Page(5)) {
		t.Errorf("IsCurrent(5) = false; want true")
	}
	if pager.IsCurrent(Page(4)) {
		t.Errorf("IsCurrent(4) = true; want false")
	}
	if got := pager.Next().String(); got != "6" {
		t.Errorf("Next() = %v; want 6", got)
	}
}

func TestGetRange(t *testing.T) {
	tests := []struct {
		name      string
		paginator *Paginator
		want      []Page
	}{
		{
			"Very few pages",
			NewPaginator(1, 40),
			[]Page{1, 2, 3},
		},
		{
			"Few pages, on the limit",
			NewPaginator(1, 20),
			[]Page{1, 2, 3, 4, 5},
		},
		{
			"Few pages, pass the limit",
			NewPaginator(1, 10),
			[]Page{1, 2, 3, -1, 10},
		},
		{
			"Many pages, current Page in middle",
			NewPaginator(10, 5),
			[]Page{1, -1, 8, 9, 10, 11, 12, -1, 20},
		},
		{
			"Many pages, current Page near start",
			NewPaginator(3, 4),
			[]Page{1, 2, 3, 4, 5, -1, 25},
		},
		{
			"Many pages, current Page near end",
			NewPaginator(23, 4),
			[]Page{1, -1, 21, 22, 23, 24, 25},
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.paginator.SetCount(100)
			if err != nil {
				t.Fatalf("SetCount() error = %v", err)
			}
			got := tt.paginator.GetRange()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRange() = %v; want %v", got, tt.want)
			}
		})
	}
}
