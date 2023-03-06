package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"trikliq-airport-finder/pkg/pdf"
	"unicode"
)

func main() {
	//server.Start()

	raw, err := ioutil.ReadFile("data/tickets/singaporeAirlines.pdf")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	txt, _ := pdf.PdfToTxt(raw)
	fmt.Println(txt)
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
	fmt.Println(newTxt)

	raw, _ = ioutil.ReadFile("data/iata.json")
	iata := make(map[string]map[string]string)

	json.Unmarshal(raw, &iata)

	candidates := make([]string, 0)
	for _, split := range splits {

		if _, found := iata[split]; found {
			candidates = append(candidates, split)
		}
	}

	fmt.Println(candidates)

	raw, _ = ioutil.ReadFile("data/moneycode.json")
	moneyCode := make(map[string]map[string]string)

	json.Unmarshal(raw, &moneyCode)

	// for _, candidate := range candidates {

	// }

	trie := pdf.TrieData()
	//Passing the words in the trie
	raw, _ = ioutil.ReadFile("data/city.json")

	word := make([]string, 0)

	json.Unmarshal(raw, &word)

	for _, wr := range word {
		trie.Insert(wr)
	}

	words := strings.Split(txt, " ")
	//initializing search for the words
	for _, wr := range words {
		found := trie.Search(wr)
		if found == 1 {
			fmt.Printf("\"%s\"Word found in trie\n", string(wr))
		}
	}
}
