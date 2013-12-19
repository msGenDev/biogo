package bwt

import "code.google.com/p/biogo/seq/linear"
import "code.google.com/p/biogo/alphabet"

type Index struct {
	alphabet alphabet.Alphabet
	BWT      []alphabet.Letter
	sa       []int
	c        [256]int
}

func New(seq *linear.Seq) *Index {
	length := seq.Len()
	suffixArray := generateSuffixArray(seq.Seq)

	index := Index{
		sa:       suffixArray,
		alphabet: seq.Alphabet(),
	}

	index.BWT = make(alphabet.Letters, length)
	for i := 0; i < length; i++ {
		dataIndex := (suffixArray[i] - 1 + length) % length
		index.BWT[i] = seq.Seq[dataIndex]
	}

	var mem alphabet.Letter

	for i, d := range suffixArray {
		currentLetter := seq.Seq[d]
		if mem != currentLetter {
			index.c[currentLetter] = i
		}
		mem = currentLetter
	}

	return &index
}

func (index *Index) SearchForBytesBasic(pattern []byte) []int {
	searchLetters := alphabet.BytesToLetters(pattern)
	s := 1
	e := len(index.BWT)

	for i := len(pattern) - 1; i >= 0; i-- {
		currChar := searchLetters[i]
		s = index.c[currChar] + index.rank(s-1, currChar) + 1
		e = index.c[currChar] + index.rank(e, currChar)
		if e <= s {
			return []int{}
		}
	}

	results := make([]int, e-s+1)
	for i := 0; i < len(results); i++ {
		results[i] = index.sa[s+i-1]
	}

	return results
}

// TODO: This is ridiculously slow.
// Implement rank and count using wavelet trees and RRR for easy speedup.
func (index *Index) rank(ranking int, c alphabet.Letter) int {
	sum := 0
	for i := 0; i < ranking; i++ {
		if index.BWT[i] == c {
			sum++
		}
	}
	return sum
}
