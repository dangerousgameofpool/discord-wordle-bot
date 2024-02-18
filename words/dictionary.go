package dictionary

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

// Dictionary represents a Wordle game's internal
// dictionary, used to check the validity of user
// guesses and select a Wordle puzzles's answer.
type Dictionary struct {
	wordlist map[string]int
	check    []string
}

// NewDictionary returns a new Dictionary struct, with
// wordlist and check fields containing words filtered
// to length l.
func NewDictionary(l int) Dictionary {
	d := Dictionary{
		wordlist: map[string]int{},
		check:    []string{},
	}

	f, err := os.Open("words/combined_wordlist.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		if len(scanner.Text()) == l {
			d.wordlist[scanner.Text()] = i
			d.check = append(d.check, scanner.Text())
			i++
		}
	}
	return d
}

// Contains checks if a Dictionary contains a given
// string s. Returns true is s is contained in the
// Dictionary, and false if it is not.
func (d Dictionary) Contains(s string) bool {
	s = strings.ToLower(s)
	_, ok := d.wordlist[s]
	return ok
}

// RandomWord selects a random string from a
// Dictionary's check field and returns it.
func (d Dictionary) RandomWord() string {
	r := rand.Intn(len(d.check))
	return d.check[r]
}
