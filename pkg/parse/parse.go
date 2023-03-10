package parse

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"trikliq-airport-finder/pkg/pdf"
	"trikliq-airport-finder/pkg/transform"
	"unicode"

	"go.uber.org/zap"
)

func Parse(raw []byte, log *zap.Logger) (finalized map[string]any) {

	txt, _ := pdf.PdfToTxt(raw)
	newTxt := ""

	for _, ch := range txt {
		if unicode.IsDigit(ch) || unicode.IsLetter(ch) {
			newTxt += string(ch)
		}
		if unicode.IsSpace(ch) {
			newTxt += " "
		}
	}

	splits := strings.Split(newTxt, " ")

	raw, _ = ioutil.ReadFile("data/iata.json")
	iata := make(map[string]map[string]string)

	json.Unmarshal(raw, &iata)

	candidates := make([]string, 0)
	for _, split := range splits {

		if _, found := iata[split]; found {
			candidates = append(candidates, split)
		}
	}

	trie := pdf.TrieData()
	//Passing the words in the trie
	raw, _ = ioutil.ReadFile("data/cityLower.json")
	cities := make([]string, 0)

	json.Unmarshal(raw, &cities)

	for _, city := range cities {
		city = strings.Replace(city, "-", " ", -1)
		trie.Insert(city)
	}

	txt = strings.Replace(txt, "\t", " ", -1)
	txt = strings.Replace(txt, "\n", " ", -1)
	txt = strings.Replace(txt, "\f", " ", -1)

	words := strings.Split(txt, " ")
	foundCities := make([]string, 0)

	//initializing search for the words
	for _, wr := range words {
		if wr == "" {
			continue
		}

		wr = strings.ToLower(wr)
		found := trie.Search(wr)
		if found == 1 {
			if !transform.InSlice(wr, foundCities) {
				foundCities = append(foundCities, wr)
			}
		}
	}

	finalCandidates := make([]string, 0)
	log.Debug("found candidates",
		zap.Strings("candidates", finalCandidates),
	)

	for _, candidate := range candidates {

		if iata[candidate] == nil {
			continue
		}

		city := iata[candidate]["city"]

		if city == "" {
			continue
		}

		find := false
		for _, found := range foundCities {
			if strings.Contains(strings.ToLower(city), found) {
				find = true
				break
			}
		}
		if find {
			finalCandidates = append(finalCandidates, candidate)
		}
	}

	finalized = Finalize(finalCandidates, iata)

	return
}
