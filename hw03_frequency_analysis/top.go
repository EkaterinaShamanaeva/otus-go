package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(inputString string) []string {
	st := strings.Fields(inputString)
	// создание словаря слов
	dict := make(map[string]int)
	for _, val := range st {
		dict[val] = 0
	}
	// подсчет количества вхождений
	for i := 0; i < len(st); i++ {
		for key := range dict {
			if st[i] == key {
				dict[key]++
			}
		}
	}
	// сортировка по количеству повторений и алфавиту
	type words struct {
		word  string
		count int
	}
	lstWords := make([]words, 0, len(dict))
	for k, v := range dict {
		lstWords = append(lstWords, words{k, v})
	}
	sort.Slice(lstWords, func(i, j int) bool {
		if lstWords[i].count == lstWords[j].count {
			return lstWords[i].word < lstWords[j].word
		}
		return lstWords[i].count > lstWords[j].count
	})
	// создание слайса результатов
	result := make([]string, 0, 10)
	if len(lstWords) >= 10 {
		for i := 0; i < 10; i++ {
			result = append(result, lstWords[i].word)
		}
	} else {
		for _, val := range lstWords {
			result = append(result, val.word)
		}
	}
	return result
}
