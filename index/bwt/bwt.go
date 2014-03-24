package bwt

import (
	"code.google.com/p/biogo/alphabet"
	"code.google.com/p/biogo/seq/linear"
	"github.com/robsyme/wavelettree"
)

type Index struct {
	alphabet alphabet.Alphabet
	BWT      []byte
	sa       []int
	c        [256]int
}

type IndexFaster struct {
	BWT   []byte
	sa    []int
	c     [256]uint
	wtree *wavelettree.WaveletTree
}

func New(seq *linear.Seq) *Index {
	length := seq.Len()
	suffixArray := generateSuffixArray(seq.Seq)

	index := Index{
		sa:       suffixArray,
		alphabet: seq.Alphabet(),
	}

	index.BWT = make([]byte, length)
	for i := 0; i < length; i++ {
		dataIndex := (suffixArray[i] - 1 + length) % length
		index.BWT[i] = byte(seq.Seq[dataIndex])
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

func NewWithWaveletTree(seq *linear.Seq) *IndexFaster {
	length := seq.Len()
	suffixArray := generateSuffixArray(seq.Seq)

	index := IndexFaster{
		sa: suffixArray,
	}

	index.BWT = make([]byte, length)
	for i := 0; i < length; i++ {
		dataIndex := (suffixArray[i] - 1 + length) % length
		index.BWT[i] = byte(seq.Seq[dataIndex])
	}

	var mem alphabet.Letter

	for i, d := range suffixArray {
		currentLetter := seq.Seq[d]
		if mem != currentLetter {
			index.c[currentLetter] = uint(i)
		}
		mem = currentLetter
	}

	index.wtree = wavelettree.New(index.BWT)
	return &index
}

func (index *Index) SearchForBytes(pattern []byte) []int {
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
		if index.BWT[i] == byte(c) {
			sum++
		}
	}
	return sum
}

func (index *IndexFaster) SearchForBytes(pattern []byte) []int {
	s := uint(1)
	e := uint(len(index.BWT))

	for i := len(pattern) - 1; i >= 0; i-- {
		currChar := pattern[i]
		s = index.c[currChar] + index.wtree.Rank(s-1, currChar) + 1
		e = index.c[currChar] + index.wtree.Rank(e, currChar)
		if e <= s {
			return []int{}
		}
	}

	results := make([]int, e-s+1)
	for i := uint(0); i < uint(len(results)); i++ {
		results[i] = index.sa[s+i-1]
	}

	return results
}

func (index *Index) rankFaster(ranking int, c alphabet.Letter) int {
	sum := 0
	for i := 0; i < ranking; i++ {
		if index.BWT[i] == byte(c) {
			sum++
		}
	}
	return sum
}
