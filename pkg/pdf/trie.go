package pdf

// Declaring trie_Node  for creating node in a trie
type trie_Node struct {
	//assigning limit of 26 for child nodes
	childrens [26]*trie_Node
	//declaring a bool variable to check the word end.
	wordEnds bool
}

// Initializing the root of the trie
type Trie struct {
	root *trie_Node
}

// inititlaizing a new trie
func TrieData() *Trie {
	t := new(Trie)
	t.root = new(trie_Node)
	return t
}

// Passing words to trie
func (t *Trie) Insert(word string) {
	current := t.root
	for _, wr := range word {

		index := wr - 'a'
		if wr < 97 || wr > 122 {
			continue
		}

		if current.childrens[index] == nil {
			current.childrens[index] = new(trie_Node)
		}
		current = current.childrens[index]
	}
	current.wordEnds = true
}

// Initializing the search for word in node
func (t *Trie) Search(word string) int {
	oneMatch := false
	current := t.root
	for _, wr := range word {
		if wr < 97 || wr > 122 {
			return 0
		}
		index := wr - 'a'

		if current.childrens[index] == nil {
			return 0
		}

		oneMatch = true
		current = current.childrens[index]
	}

	if current.wordEnds && oneMatch {
		return 1
	}

	return 0
}
