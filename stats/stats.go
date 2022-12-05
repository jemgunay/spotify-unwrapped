package stats

import (
	"sort"
	"time"

	"github.com/jemgunay/spotify-unwrapped/spotify"
)

// Detail represents a mapping of a track to its stat value.
type Detail struct {
	id         string
	Name       string  `json:"name,omitempty"`
	CoverImage string  `json:"cover_image,omitempty"`
	SpotifyURL string  `json:"spotify_url,omitempty"`
	Value      float64 `json:"value"`
	Date       string  `json:"date,omitempty"`
}

func (d *Detail) Set(name, image, spotifyURL string) {
	d.Name = name
	d.CoverImage = image
	d.SpotifyURL = spotifyURL
}

// DateYear determines the year from the date value.
func (d *Detail) DateYear() int {
	return time.Unix(int64(d.Value), 0).Year()
}

// Group is used to calculate the min, max and average values for a dataset.
type Group struct {
	Min   Detail `json:"min"`
	Max   Detail `json:"max"`
	sum   float64
	count float64
	Mean  Detail `json:"avg"`
}

// Push pushes a value and its key into the Group. Call Calc to finalise the Group statistics.
func (g *Group) Push(id string, val float64) {
	g.sum += val
	g.count++

	d := Detail{id: id, Value: val}
	switch {
	case g.Min.id == "":
		g.Min, g.Max = d, d
	case val > g.Max.Value:
		g.Max = d
	case val < g.Min.Value:
		g.Min = d
	}
}

// GroupCalcOpt defines a Calc option.
type GroupCalcOpt func(*Group)

// WithMultiplier multiplies the Group output values, e.g. can be used to convert decimal values to percentages.
func WithMultiplier(multiplier float64) GroupCalcOpt {
	return func(group *Group) {
		group.Min.Value *= multiplier
		group.Max.Value *= multiplier
		group.Mean.Value *= multiplier
	}
}

// Calc calculates the final statistics for the Group; to be called once all values have bene Pushed.
func (g *Group) Calc(lookup map[string]spotify.TrackDetails, opts ...GroupCalcOpt) {
	minTrack := lookup[g.Min.id]
	maxTrack := lookup[g.Max.id]
	g.Min.Set(minTrack.GetTrackString(), minTrack.Album.Images.First(), minTrack.ExternalURLs.Spotify)
	g.Max.Set(maxTrack.GetTrackString(), maxTrack.Album.Images.First(), maxTrack.ExternalURLs.Spotify)
	if g.count > 0 {
		g.Mean.Value = g.sum / g.count
	}

	for _, opt := range opts {
		opt(g)
	}
}

// CalcDate is Calc but sets the Group Date fields to a formatted timestamp.
func (g *Group) CalcDate(lookup map[string]spotify.TrackDetails) {
	g.Calc(lookup)

	g.Min.Date = unixToDate(g.Min.Value)
	g.Max.Date = unixToDate(g.Max.Value)
	if g.count > 0 {
		g.Mean.Date = unixToDate(g.Mean.Value)
	}
}

func unixToDate(val float64) string {
	return time.Unix(int64(val), 0).Format("02/01/2006")
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

// MappingOpt defines operations to be performed during Mapping.OrderedLabelsAndValues calls.
type MappingOpt func(*OrderedKVPair)

// SortBy defines the OrderedKVPair sort types.
type SortBy int

const (
	// SortKey sorts the OrderedKVPair by key.
	SortKey SortBy = iota
	// SortValue sorts the OrderedKVPair by value.
	SortValue
)

// WithSort sorts OrderedKVPair by the provided sort type.
func WithSort(sortBy SortBy, sortDesc bool) MappingOpt {
	return func(pair *OrderedKVPair) {
		pair.sortBy = sortBy
		pair.sortDesc = sortDesc
		sort.Sort(pair)
	}
}

// WithTruncate truncates the OrderedKVPair to the specified length.
func WithTruncate(size int) MappingOpt {
	return func(pair *OrderedKVPair) {
		if size > len(pair.Keys) {
			size = len(pair.Keys)
		}
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
		if p.Keys[i] == p.Keys[j] {
			return p.Values[i] < p.Values[j]
		}
		return p.Keys[i] < p.Keys[j]
	}
	if p.Values[i] == p.Values[j] {
		return p.Keys[i] < p.Keys[j]
	}
	return p.Values[i] < p.Values[j]
}

// Swap swaps two key/value pairs.
func (p *OrderedKVPair) Swap(i int, j int) {
	p.Keys[i], p.Keys[j] = p.Keys[j], p.Keys[i]
	p.Values[i], p.Values[j] = p.Values[j], p.Values[i]
}
