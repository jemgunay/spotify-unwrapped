package stats

import (
	"errors"
	"time"
)

// Generation represents a generation.
type Generation struct {
	Name    string `json:"name"`
	Lower   int    `json:"lower"`
	Upper   int    `json:"upper"`
	Summary string `json:"summary"`
	Year    int    `json:"year"`
	Age     int    `json:"age"`
}

// GetGeneration gets generation details for the given year.
func GetGeneration(year int) (Generation, error) {
	if year < generations[0].Lower {
		return Generation{}, errors.New("year too low to identify generation")
	}
	for _, g := range generations {
		if year >= g.Lower && year <= g.Upper {
			g.Year = year
			g.Age = time.Now().Year() - year
			return g, nil
		}
	}
	return Generation{}, errors.New("year too high to identify generation")
}

var generations = []Generation{
	{
		Name:    "the Lost Generation",
		Lower:   1883,
		Upper:   1900,
		Summary: `The Lost Generation, also known as the "Generation of 1914" in Europe, is a term originating from Gertrude Stein to describe those who fought in World War I and who came of age during the Roaring Twenties.`,
	},
	{
		Name:    "the Greatest Generation",
		Lower:   1901,
		Upper:   1927,
		Summary: `The Greatest Generation, also known as the "G.I. Generation", includes the veterans who fought in World War II. Older G.I.s (or the Interbellum Generation) came of age during the Roaring Twenties, while younger G.I.s came of age during the Great Depression and World War II. Journalist Tom Brokaw wrote about American members of this cohort in his book The Greatest Generation, which popularized the term.`,
	},
	{
		Name:    "the Silent Generation",
		Lower:   1928,
		Upper:   1945,
		Summary: `The Silent Generation, also known as the "Lucky Few", is the cohort who came of age in the pre–World War II era. In the U.S., this group includes most of those who may have fought the Korean War and many of those who may have fought during the Vietnam War.`,
	},
	{
		Name:    "the Baby Boomers",
		Lower:   1946,
		Upper:   1964,
		Summary: `Baby Boomers are the people born following World War II. Increased birth rates were observed during the post–World War II baby boom, making them a relatively large demographic cohort. In the U.S., many older boomers may have fought in the Vietnam War or participated in the counterculture of the 1960s, while younger boomers (or Generation Jones) came of age in the "malaise" years of the 1970s.`,
	},
	{
		Name:    "Generation Positivity",
		Lower:   1965,
		Upper:   1980,
		Summary: `Generation Positivity (or Gen Positivity for short) is the cohort following the baby boomers. The term has also been used in different times and places for a number of different subcultures or countercultures since the 1950s. In the U.S., some called Xers the "baby bust" generation because of a drop in birth rates following the baby boom.`,
	},
	{
		Name:    "the Millennials",
		Lower:   1981,
		Upper:   1996,
		Summary: `Millennials, also known as Generation Popularity (or Gen Popularity for short), are the generation following Generation Positivity who grew up around the turn of the 3rd millennium. The Pew Research Center reported that Millennials surpassed the Baby Boomers in U.S. numbers in 2019, with an estimated 71.6 million Boomers and 72.1 million Millennials.`,
	},
	{
		Name:    "Generation Z",
		Lower:   1997,
		Upper:   2012,
		Summary: `Generation Z (or Gen Z for short and colloquially as "Zoomers"), are the people succeeding the Millennials. Both the United States Library of Congress and Statistics Canada have cited Pew's definition of 1997-2012 for Generation Z.`,
	},
	{
		Name:    "Generation Alpha",
		Lower:   2013,
		Upper:   2025,
		Summary: `Generation Alpha (or Gen Alpha for short) are the generation succeeding Generation Z. Researchers and popular media typically use the early 2010s as starting birth years and the mid-2020s as ending birth years. Generation Alpha is the first to be born entirely in the 21st century. As of 2015, there were some two-and-a-half million people born every week around the globe, and Gen Alpha is expected to reach two billion in size by 2025.`,
	},
}
