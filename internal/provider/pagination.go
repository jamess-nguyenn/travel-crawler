package provider

type Pagination struct {
	// Total matching records across all pages.
	// nil means unknown.
	Total *int

	// Whether another page exists.
	// nil means unknown.
	HasMore *bool

	// Cursor for fetching the next page.
	// nil means unavailable.
	NextCursor *string
}

// TotalValue returns the total number of matching records.
// The second return value indicates whether Total is known.
func (p *Pagination) TotalValue() (int, bool) {
	if p.Total == nil {
		return 0, false
	}

	return *p.Total, true
}

// HasMoreValue returns whether another page exists.
// The second return value indicates whether HasMore is known.
func (p *Pagination) HasMoreValue() (bool, bool) {
	if p.HasMore == nil {
		return false, false
	}

	return *p.HasMore, true
}

// NextCursorValue returns the cursor for the next page.
// The second return value indicates whether a non-empty cursor exists.
func (p *Pagination) NextCursorValue() (string, bool) {
	if p.NextCursor == nil || *p.NextCursor == "" {
		return "", false
	}

	return *p.NextCursor, true
}

// HasNextPage reports whether another page exists and whether
// that information is known.
//
// HasMore takes precedence when explicitly provided.
// Otherwise, a non-empty cursor indicates that another page exists.
func (p *Pagination) HasNextPage() (bool, bool) {
	if p.HasMore != nil {
		return *p.HasMore, true
	}

	if _, ok := p.NextCursorValue(); ok {
		return true, true
	}

	return false, false
}

// SetTotal sets the total number of matching records.
func (p *Pagination) SetTotal(total int) {
	p.Total = &total
}

// SetHasMore sets whether another page exists.
func (p *Pagination) SetHasMore(hasMore bool) {
	p.HasMore = &hasMore
}

// SetNextCursor sets the cursor for the next page.
func (p *Pagination) SetNextCursor(cursor string) {
	p.NextCursor = &cursor
}
