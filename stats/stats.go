package stats

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/jemgunay/spotify-unwrapped/spotify"
)

// Detail represents a mapping of a track to its stat value and other track metadata.
type Detail struct {
	id         string
	Name       string `json:"name,omitempty"`
	CoverImage string `json:"cover_image,omitempty"`
	SpotifyURL string `json:"spotify_url,omitempty"`
	// value stores the raw float data, ValueOut represents the processed output value, e.g. in date/duration formats.
	value    float64
	ValueOut any `json:"value"`
}

// Set sets track metadata for the given Detail.
func (d *Detail) Set(name, image, spotifyURL string) {
	d.Name = name
	d.CoverImage = image
	d.SpotifyURL = spotifyURL
}

// DateYear determines the year from the date value.
func (d *Detail) DateYear() int {
	return time.Unix(int64(d.value), 0).Year()
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

	d := Detail{id: id, value: val}
	switch {
	case g.Min.id == "":
		g.Min, g.Max = d, d
	case val > g.Max.value:
		g.Max = d
	case val < g.Min.value:
		g.Min = d
	}
}

// GroupCalcOpt defines a Calc option.
type GroupCalcOpt func(*Group)

// WithMultiplier multiplies the Group output values, e.g. can be used to convert decimal values to percentages.
func WithMultiplier(multiplier float64) GroupCalcOpt {
	return func(group *Group) {
		group.Min.value *= multiplier
		group.Max.value *= multiplier
		group.Mean.value *= multiplier
	}
}

// ToDateString sets the output value to the float value processed into a date string.
func ToDateString() GroupCalcOpt {
	return func(group *Group) {
		group.Min.ValueOut = unixToDate(group.Min.value)
		group.Max.ValueOut = unixToDate(group.Max.value)
		if group.count > 0 {
			group.Mean.ValueOut = unixToDate(group.Mean.value)
		}
	}
}

func unixToDate(unixDate float64) string {
	return time.Unix(int64(unixDate), 0).Format("02/01/2006")
}

// ToDurationString sets the output value to the float value processed into a duration string.
func ToDurationString() GroupCalcOpt {
	return func(group *Group) {
		group.Min.ValueOut = secondsToDuration(group.Min.value)
		group.Max.ValueOut = secondsToDuration(group.Max.value)
		if group.count > 0 {
			group.Mean.ValueOut = secondsToDuration(group.Mean.value)
		}
	}
}

var durationReplacer = strings.NewReplacer(
	"h", "h ",
	"m", "m ",
	"s", "s ",
)

func secondsToDuration(seconds float64) string {
	dur := time.Duration(seconds) * time.Millisecond
	// Duration.String() doesn't provide spaces - add them in
	dur = dur.Round(time.Second)
	repl := durationReplacer.Replace(dur.String())
	return strings.TrimRight(repl, " ")
}

// Calc calculates the final statistics for the Group; to be called once all values have bene Pushed.
func (g *Group) Calc(lookup map[string]spotify.TrackDetails, opts ...GroupCalcOpt) {
	minTrack, maxTrack := lookup[g.Min.id], lookup[g.Max.id]
	g.Min.Set(minTrack.GetTrackString(), minTrack.Album.Images.First(), minTrack.ExternalURLs.Spotify)
	g.Max.Set(maxTrack.GetTrackString(), maxTrack.Album.Images.First(), maxTrack.ExternalURLs.Spotify)
	if g.count > 0 {
		g.Mean.value = g.sum / g.count
	}

	g.Min.ValueOut = nil
	for _, opt := range opts {
		opt(g)
	}

	// some opts set ValueOut, some don't. If not set by an opt, set to the base float value.
	if g.Min.ValueOut == nil {
		g.Min.ValueOut = math.Round(g.Min.value)
		g.Max.ValueOut = math.Round(g.Max.value)
		g.Mean.ValueOut = math.Round(g.Mean.value)
	}
}

// Mapping maps a key to a count of the occurrences of that key.
type Mapping map[string]int

// NewMapping returns an initialised Mapping with an initial capacity.
func NewMapping(capacity int, defaultKeys ...string) Mapping {
	m := make(map[string]int, capacity)
	for _, k := range defaultKeys {
		m[k] = 0
	}
	return m
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
	// SortPitchKey sorts by key, but using the musical pitch key notation order.
	SortPitchKey
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
		// reverse ordering
		i, j = j, i
	}

	// sort by key
	if p.sortBy == SortKey {
		if p.Keys[i] == p.Keys[j] {
			return p.Values[i] < p.Values[j]
		}
		return p.Keys[i] < p.Keys[j]
	}

	if p.sortBy == SortPitchKey {
		// drive the comparison by converting the standard pitch notation back to Spotify's ordered integer key values
		return pitchKeyToIntMappings[p.Keys[i]] < pitchKeyToIntMappings[p.Keys[j]]
	}

	// sort by value (default fallback)
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

var (
	// PitchKeys is the set of musical pitch keys. Spotify returns them as zero-indexed integers which map to the
	// pitches using standard Pitch Class notation. E.g. 0 = C, 1 = C♯/D♭, 2 = D, etc. If no key was detected, the
	// value is -1
	PitchKeys = []string{
		"C",
		"C♯/D♭",
		"D",
		"D♯/E♭",
		"E",
		"F",
		"F♯/G♭",
		"G",
		"G♯/A♭",
		"A",
		"A♯/B♭",
		"B",
	}
	intToPitchKeyMappings = make(map[int]string, len(PitchKeys))
	pitchKeyToIntMappings = make(map[string]int, len(PitchKeys))
)

func init() {
	for k, v := range PitchKeys {
		intToPitchKeyMappings[k] = v
		pitchKeyToIntMappings[v] = k
	}
}

// SpotifyKeyToPitchKey maps Spotify's 0-11 integer pitch key values to the standard pitch class notation strings.
func SpotifyKeyToPitchKey(i int) string {
	return intToPitchKeyMappings[i]
}
