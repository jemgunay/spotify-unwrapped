package stats

import (
	"sort"
	"time"

	"github.com/jemgunay/spotify-unwrapped/spotify"
)

type detail struct {
	id    string
	Name  string  `json:"name,omitempty"`
	Value float64 `json:"value,omitempty"`
	Date  string  `json:"date,omitempty"`
}

// DateYear determines the year from the date value.
func (d detail) DateYear() int {
	return time.Unix(int64(d.Value), 0).Year()
}

// Group is used to calculate the min, max and average values for a dataset.
type Group struct {
	Min   detail `json:"min"`
	Max   detail `json:"max"`
	sum   float64
	count float64
	Mean  detail `json:"avg"`
}

// Push pushes a value and its key into the Group. Call Calc to finalise the Group statistics.
func (g *Group) Push(id string, val float64) {
	g.sum += val
	g.count++

	d := detail{id: id, Value: val}
	switch {
	case g.Min.id == "":
		g.Min, g.Max = d, d
	case val > g.Max.Value:
		g.Max = d
	case val < g.Min.Value:
		g.Min = d
	}
}

// Calc calculates the final statistics for the Group; to be called once all values have bene Pushed.
func (g *Group) Calc(lookup map[string]spotify.TrackDetails) {
	minTrack := lookup[g.Min.id]
	g.Min.Name = minTrack.GetTrackString()
	maxTrack := lookup[g.Max.id]
	g.Max.Name = maxTrack.GetTrackString()
	if g.count > 0 {
		g.Mean.Value = g.sum / g.count
	}
}

func (g *Group) CalcDate(lookup map[string]spotify.TrackDetails) {
	minTrack := lookup[g.Min.id]
	g.Min.Name = minTrack.GetTrackString()
	g.Min.Date = unixToDate(g.Min.Value)
	maxTrack := lookup[g.Max.id]
	g.Max.Name = maxTrack.GetTrackString()
	g.Max.Date = unixToDate(g.Max.Value)
	if g.count > 0 {
		g.Mean.Value = g.sum / g.count
		g.Mean.Date = unixToDate(g.Mean.Value)
	}
}

func unixToDate(val float64) string {
	return time.Unix(int64(val), 0).Format("02-01-2006")
}

// Mapping maps a key to a count of the occurrences of that key.
type Mapping map[string]int

// NewMapping returns an initialised Mapping with an initial capacity.
func NewMapping(capacity int) Mapping {
	return make(map[string]int, capacity)
}

// Push increments the count for the given key.
func (m Mapping) Push(key string) {
	m[key] = m[key] + 1
}

type MappingOpt func(*OrderedKVPair)

type SortBy int

const (
	SortKey SortBy = iota
	SortValue
)

func WithSort(sortBy SortBy, sortDesc bool) MappingOpt {
	return func(pair *OrderedKVPair) {
		pair.sortBy = sortBy
		pair.sortDesc = sortDesc
		sort.Sort(pair)
	}
}

func WithTruncate(size int) MappingOpt {
	return func(pair *OrderedKVPair) {
		pair.Keys = pair.Keys[:size]
		pair.Values = pair.Values[:size]
	}
}

// OrderedLabelsAndValues converts a Mapping to an OrderedKVPair for use with ChartJS.
func (m Mapping) OrderedLabelsAndValues(opts ...MappingOpt) *OrderedKVPair {
	pair := &OrderedKVPair{
		Keys:   make([]string, 0, len(m)),
		Values: make([]int, 0, len(m)),
	}
	for k, v := range m {
		pair.Keys = append(pair.Keys, k)
		pair.Values = append(pair.Values, v)
	}
	for _, opt := range opts {
		opt(pair)
	}
	return pair
}

// OrderedKVPair provides two lists of keys and their corresponding values, ordered by key. This format is required by
// ChartJS. It implements sort.Interface.
type OrderedKVPair struct {
	Keys     []string `json:"keys"`
	Values   []int    `json:"values"`
	sortBy   SortBy
	sortDesc bool
}

var _ sort.Interface = (*OrderedKVPair)(nil)

// Len returns the number of keys.
func (p *OrderedKVPair) Len() int {
	return len(p.Keys)
}

// Less orders the provided elements.
func (p *OrderedKVPair) Less(i int, j int) bool {
	if p.sortDesc {
		i, j = j, i
	}
	if p.sortBy == SortKey {
		return p.Keys[i] < p.Keys[j]
	}
	return p.Values[i] < p.Values[j]
}

// Swap swaps two key/value pairs.
func (p *OrderedKVPair) Swap(i int, j int) {
	p.Keys[i], p.Keys[j] = p.Keys[j], p.Keys[i]
	p.Values[i], p.Values[j] = p.Values[j], p.Values[i]
}
