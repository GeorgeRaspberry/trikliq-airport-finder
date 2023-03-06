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

	// var (
	// 	filePath = "file.json"
	// 	airports = make(map[string]map[string]string)
	// )

	// file, err := os.Open(filePath)
	// if err != nil {
	// 	fmt.Printf("failed to open a file")
	// }

	// bytes, _ := ioutil.ReadAll(file)

	// _ = json.Unmarshal(bytes, &airports)

	// icao := make(map[string]interface{})
	// iata := make(map[string]interface{})

	// for _, value := range airports {

	// 	if value["iata"] != "" {
	// 		iata[value["iata"]] = value
	// 	}

	// 	if value["icao"] != "" {
	// 		icao[value["icao"]] = value
	// 	}
	// }
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

	// airportsBytes, _ := json.Marshal(iata)
	// ioutil.WriteFile("iata.json", airportsBytes, 0744)

	// airportsBytes, _ = json.Marshal(icao)
	// ioutil.WriteFile("icao.json", airportsBytes, 0744)

	// fmt.Printf("JSON data written to file")

	trie := pdf.TrieData()
	//Passing the words in the trie
	word := []string{"aqua", "jack", "card", "care"}
	for _, wr := range word {
		trie.Insert(wr)
	}
	//initializing search for the words
	words_Search := []string{"aqua", "jack", "card", "care", "cat", "dog", "can"}
	for _, wr := range words_Search {
		found := trie.Search(wr)
		if found == 1 {
			fmt.Printf("\"%s\"Word found in trie\n", wr)
		} else {
			fmt.Printf(" \"%s\" Word not found in trie\n", wr)
		}
	}
}
