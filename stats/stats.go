package stats

import (
	"sort"

	"github.com/jemgunay/spotify-unwrapped/spotify"
)

type detail struct {
	id    string
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// Group is used to calculate the min, max and average values for a dataset.
type Group struct {
	Min   detail `json:"min"`
	Max   detail `json:"max"`
	sum   float64
	count float64
	Mean  float64 `json:"avg"`
}

// Push pushes a value and its key into the Group. Call Calc to finalise the Group statistics.
func (s *Group) Push(id string, val float64) {
	s.sum += val
	s.count++

	switch {
	case s.Min.id == "":
		s.Min = detail{id: id, Value: val}
		s.Max = detail{id: id, Value: val}
	case val > s.Max.Value:
		s.Max = detail{id: id, Value: val}
	case val < s.Min.Value:
		s.Min = detail{id: id, Value: val}
	}
}

// Calc calculates the final statistics for the Group; to be called once all values have bene Pushed.
func (s *Group) Calc(lookup map[string]spotify.TrackDetails) {
	minTrack := lookup[s.Min.id]
	s.Min.Name = minTrack.GetTrackString()
	maxTrack := lookup[s.Max.id]
	s.Max.Name = maxTrack.GetTrackString()
	s.Mean = s.sum / s.count
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

// OrderedLabelsAndValues converts a Mapping to an OrderedKVPair for use with ChartJS.
func (m Mapping) OrderedLabelsAndValues() OrderedKVPair {
	pair := OrderedKVPair{
		Keys:   make([]string, 0, len(m)),
		Values: make([]int, 0, len(m)),
	}
	for k, v := range m {
		pair.Keys = append(pair.Keys, k)
		pair.Values = append(pair.Values, v)
	}
	sort.Sort(&pair)
	return pair
}

// OrderedKVPair provides two lists of keys and their corresponding values, ordered by key. This format is required by
// ChartJS. It implements sort.Interface.
type OrderedKVPair struct {
	Keys   []string `json:"keys"`
	Values []int    `json:"values"`
}

var _ sort.Interface = (*OrderedKVPair)(nil)

// Len returns the number of keys.
func (p *OrderedKVPair) Len() int {
	return len(p.Keys)
}

// Less orders the provided elements.
func (p *OrderedKVPair) Less(i int, j int) bool {
	return p.Keys[i] < p.Keys[j]
}

// Swap swaps two key/value pairs.
func (p *OrderedKVPair) Swap(i int, j int) {
	p.Keys[i], p.Keys[j] = p.Keys[j], p.Keys[i]
	p.Values[i], p.Values[j] = p.Values[j], p.Values[i]
}
