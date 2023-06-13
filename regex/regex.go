package regex

import (
	"regexp"
	"strings"
	"sync"
)

var quotes = regexp.MustCompile(`["'](\w+)["']`)
var words = regexp.MustCompile(`(\s)(\w+)(\s)?:`)
var variable = regexp.MustCompile(`(?:var|let|const)\s*(\w+)`)
var equal = regexp.MustCompile(`(\w+)(\s)?=`)
var queryParams = regexp.MustCompile(`\?(\w+)=[^&\s]+`)

func Regex(input string) []string {
	data := strings.Split(input, "\n")

	var wg sync.WaitGroup
	var mu sync.Mutex

	output := make([]string, 0, len(data))

	for _, v := range data {
		wg.Add(1)

		go func(v string) {
			defer wg.Done()

			matches := make([]string, 0, 5)

			if quotesMatches := quotes.FindAllStringSubmatch(v, -1); len(quotesMatches) > 0 {
				for _, match := range quotesMatches {
					matches = append(matches, match[1])
				}
			}
			if wordMatches := words.FindAllStringSubmatch(v, -1); len(wordMatches) > 0 {
				for _, match := range wordMatches {
					final := strings.ReplaceAll(match[0], ":", "")
					matches = append(matches, final)
				}
			}
			if variableMatches := variable.FindAllStringSubmatch(v, -1); len(variableMatches) > 0 {
				for _, match := range variableMatches {
					matches = append(matches, match[1])
				}
			}
			if equalMatches := equal.FindAllStringSubmatch(v, -1); len(equalMatches) > 0 {
				for _, match := range equalMatches {
					final := strings.ReplaceAll(match[0], "=", "")
					matches = append(matches, final)
				}
			}
			if queryParamsMatches := queryParams.FindAllStringSubmatch(v, -1); len(queryParamsMatches) > 0 {
				for _, match := range queryParamsMatches {
					matches = append(matches, match[1])
				}
			}

			mu.Lock()
			output = append(output, matches...)
			mu.Unlock()
		}(v)
	}

	wg.Wait()

	return output
}
