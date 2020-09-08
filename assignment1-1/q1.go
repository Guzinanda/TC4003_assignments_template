package cos418_hw1_1

//package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// * Find the top K most common words in a text document.
// 		path: 			location of the document
//		numWords: 		number of words to return (i.e. k)
//		charThreshold: 	character threshold for whether a token qualifies as a word, e.g.
// 						charThreshold = 5 means "apple" is a word but "pear" is not.

// * Matching is case insensitive, e.g. "Orange" and "orange" is considered the same word.
// * A word comprises alphanumeric characters only. All punctuations and other characters
//   are removed, e.g. "don't" becomes "dont".
// * You should use `checkError` to handle potential errors.

func CleanWord(dirtyWord string) string {

	// Lo hace chiquito:
	littleWord := strings.ToLower(dirtyWord)

	// Make a Regex (expresion regular) to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	cleanedWord := reg.ReplaceAllString(littleWord, "")
	return cleanedWord
}

func topWords(path string, numWords int, charThreshold int) []WordCount {
	// TODO: implement me
	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful
	// HINT: To keep only alphanumeric characters, use the regex "[^0-9a-zA-Z]+"

	file, err := os.Open(path)
	checkError(err)    // esta funcion está en common.go
	defer file.Close() // "defer" hace que esto se corra al final aunque este en esta linea
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	listofWords := make(map[string]int)

	for scanner.Scan() {
		scannedWord := CleanWord(scanner.Text())
		if len(scannedWord) >= charThreshold {
			listofWords[scannedWord]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var limitedList []WordCount
	for key, element := range listofWords {
		wordCount := WordCount{key, element}
		limitedList = append(limitedList, wordCount)
	}
	sortWordCounts(limitedList)

	return limitedList[:numWords]
}

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.
// DO NOT MODIFY THIS FUNCTION!
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
