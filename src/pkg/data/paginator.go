package data

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

const (
	paginatorEllipsis   = "…"
	paginatorOnEachSide = 2
	paginatorOnEnds     = 1
)

var ErrInvalidPage = errors.New("invalid page")

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Page int

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Paginator struct {
	onEachSide int
	onEnds     int
	NumPages   int
	Number     int
	PerPage    int
	Count      int
}

func (p Page) String() string {
	if p == -1 {
		return paginatorEllipsis
	}
	return strconv.Itoa(int(p))
}

func (p Page) IsEllipsis() bool {
	return p == -1
}

func NewPaginator(number, perPage int) *Paginator {
	return &Paginator{
		onEachSide: paginatorOnEachSide,
		onEnds:     paginatorOnEnds,
		Number:     number,
		PerPage:    perPage,
	}
}

// SetCount activates the Paginator. The Paginator won't work without it.
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (p *Paginator) SetCount(count int) error {
	p.Count = count
	return p.setNumPages()
}

func (p *Paginator) setNumPages() error {
	if p.Count == 0 {
		p.NumPages = 1
	} else {
		p.NumPages = int(math.Ceil(float64(p.Count) / float64(p.PerPage)))
	}
	switch {
	case p.Number == 0:
		p.Number = 1
	case p.Number > p.NumPages:
		return fmt.Errorf("%w: %d > %d", ErrInvalidPage, p.Number, p.NumPages)
	case p.Number < 1:
		return fmt.Errorf("%w: %d < 1", ErrInvalidPage, p.Number)
	}
	return nil
}

func (p *Paginator) getLimitOffset() (int, int) {
	offset := (p.Number - 1) * p.PerPage
	return p.PerPage, offset
}

func (p *Paginator) HasNext() bool {
	return p.Number < p.NumPages
}

func (p *Paginator) HasPrevious() bool {
	return p.Number > 1
}

func (p *Paginator) Next() Page {
	return Page(p.Number + 1)
}

func (p *Paginator) Previous() Page {
	return Page(p.Number - 1)
}

func (p *Paginator) IsCurrent(page Page) bool {
	return Page(p.Number) == page
}

// GetRange returns a 1-based range of pages with some values elided (represented by -1 / "…").
//
// [Copied from Django's paginator.py]
//
// If the Page range is larger than a given size, the whole range is not
// provided and a compact form is returned instead, e.g. for a paginator
// with 50 pages, if Page 43 were the current Page, the output would be:
//
//	1, 2, …, 40, 41, 42, 43, 44, 45, 46, …, 49, 50.
func (p *Paginator) GetRange() []Page {
	if p.NumPages <= (p.onEachSide+p.onEnds)*2 {
		return p.getSimpleRange()
	}

	var result []Page

	// Add pages at the start
	result = append(result, generateRange(1, p.onEnds)...)

	// Add pages around the current Page
	leftIndex := maxInt(p.Number-p.onEachSide, p.onEnds+1)
	rightIndex := minInt(p.Number+p.onEachSide, p.NumPages-p.onEnds)
	if leftIndex > p.onEnds+1 {
		result = append(result, -1)
	}
	result = append(result, generateRange(leftIndex, rightIndex)...)
	if rightIndex < p.NumPages-p.onEnds {
		result = append(result, -1)
	}

	// Add pages at the end
	result = append(result, generateRange(p.NumPages-p.onEnds+1, p.NumPages)...)

	return result
}

func (p *Paginator) getSimpleRange() []Page {
	return generateRange(1, p.NumPages)
}

func generateRange(start, end int) []Page {
	result := make([]Page, end-start+1)
	for i := range result {
		result[i] = Page(start + i)
	}
	return result
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
