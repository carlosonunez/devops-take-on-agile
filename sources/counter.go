package sources

import (
	"regexp"
	"strings"
)

var stopWords = map[string]int{
	"i":          0,
	"me":         0,
	"my":         0,
	"myself":     0,
	"we":         0,
	"our":        0,
	"ours":       0,
	"ourselves":  0,
	"you":        0,
	"your":       0,
	"yours":      0,
	"yourself":   0,
	"yourselves": 0,
	"he":         0,
	"him":        0,
	"his":        0,
	"himself":    0,
	"she":        0,
	"her":        0,
	"hers":       0,
	"herself":    0,
	"it":         0,
	"its":        0,
	"itself":     0,
	"they":       0,
	"them":       0,
	"their":      0,
	"theirs":     0,
	"themselves": 0,
	"what":       0,
	"which":      0,
	"who":        0,
	"whom":       0,
	"this":       0,
	"that":       0,
	"these":      0,
	"those":      0,
	"am":         0,
	"is":         0,
	"are":        0,
	"was":        0,
	"were":       0,
	"be":         0,
	"been":       0,
	"being":      0,
	"have":       0,
	"has":        0,
	"had":        0,
	"having":     0,
	"do":         0,
	"does":       0,
	"did":        0,
	"doing":      0,
	"a":          0,
	"an":         0,
	"the":        0,
	"and":        0,
	"but":        0,
	"if":         0,
	"or":         0,
	"because":    0,
	"as":         0,
	"until":      0,
	"while":      0,
	"of":         0,
	"at":         0,
	"by":         0,
	"for":        0,
	"with":       0,
	"about":      0,
	"against":    0,
	"between":    0,
	"into":       0,
	"through":    0,
	"during":     0,
	"before":     0,
	"after":      0,
	"above":      0,
	"below":      0,
	"to":         0,
	"from":       0,
	"up":         0,
	"down":       0,
	"in":         0,
	"out":        0,
	"on":         0,
	"off":        0,
	"over":       0,
	"under":      0,
	"again":      0,
	"further":    0,
	"then":       0,
	"once":       0,
	"here":       0,
	"there":      0,
	"when":       0,
	"where":      0,
	"why":        0,
	"how":        0,
	"all":        0,
	"any":        0,
	"both":       0,
	"each":       0,
	"few":        0,
	"more":       0,
	"most":       0,
	"other":      0,
	"some":       0,
	"such":       0,
	"no":         0,
	"nor":        0,
	"not":        0,
	"only":       0,
	"own":        0,
	"same":       0,
	"so":         0,
	"than":       0,
	"too":        0,
	"very":       0,
	"s":          0,
	"t":          0,
	"can":        0,
	"will":       0,
	"just":       0,
	"don":        0,
	"should":     0,
	"now":        0,
}

func countWordsInComments(cl []string) map[string]int {
	counts := map[string]int{}
	moreThanOne := map[string]int{}
	wholeWord := regexp.MustCompile("^[a-z]{2,}$")
	special := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	for _, c := range cl {
		for _, w := range strings.Split(special.ReplaceAllString(c, " "), " ") {
			w = special.ReplaceAllString(w, "")
			w = strings.ToLower(w)
			if !wholeWord.Match([]byte(w)) {
				continue
			}
			if _, found := stopWords[w]; found {
				continue
			}
			if _, found := moreThanOne[w]; found {
				if _, found := counts[w]; found {
					counts[w] = counts[w] + 1
				} else {
					counts[w] = 2
				}
			} else {
				moreThanOne[w] = 1
			}
		}
	}
	return counts
}
