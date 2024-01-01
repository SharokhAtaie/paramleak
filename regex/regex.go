package regex

import (
	"regexp"
	"strings"
	"sync"
)

var quotes = regexp.MustCompile(`["'](\w{1,50})["']`)
var words = regexp.MustCompile(`(\s)(\w{1,50})(\s)?:`)
var variable = regexp.MustCompile(`(?:var|let|const)\s*(\w{1,50})`)
var equal = regexp.MustCompile(`(\w{1,50})(\s)?=`)
var queryParams = regexp.MustCompile(`\?(\w{1,50})=[^&\s]+`)
var function = regexp.MustCompile(`\((['"]?)(\w+)(['"]?)\)`)

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
					word := strings.ReplaceAll(match[0], ":", "")
					final := strings.TrimSpace(word)
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
			if functionMatches := function.FindAllStringSubmatch(v, -1); len(functionMatches) > 0 {
				for _, match := range functionMatches {
					match[0] = strings.ReplaceAll(match[0], "(", "")
					match[0] = strings.ReplaceAll(match[0], ")", "")
					matches = append(matches, match[0])
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
