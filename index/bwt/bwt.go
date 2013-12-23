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
	wtree    *wavelettree.WaveletTree
}

type IndexFaster struct {
	alphabet alphabet.Alphabet
	BWT      []byte
	sa       []int
	c        [256]uint
	wtree    *wavelettree.WaveletTree
}

func New(seq *linear.Seq) *Index {
	data := append(seq.Seq, 0)
	length := len(data)
	suffixArray := generateSuffixArray(data)
	tree := wavelettree.New(alphabet.LettersToBytes(data))
	index := Index{
		sa:       suffixArray,
		alphabet: seq.Alphabet(),
		wtree:    tree,
	}

	index.BWT = make([]byte, length)
	for i := 0; i < length; i++ {
		dataIndex := (suffixArray[i] - 1 + length) % length
<<<<<<< HEAD
		index.BWT[i] = data[dataIndex]
=======
		index.BWT[i] = byte(seq.Seq[dataIndex])
>>>>>>> 8c03a7ded2ff0b2901294810386cdd046be11ab4
	}

	var mem alphabet.Letter

	for i, d := range suffixArray {
		currentLetter := data[d]
		if mem != currentLetter {
			index.c[currentLetter] = i
		}
		mem = currentLetter
	}

	return &index
}

<<<<<<< HEAD
func (index *Index) SeachForBytesFast(pattern []byte) []uint {
	searchLetters := alphabet.BytesToLetters(pattern)
	s := uint(1)
	e := uint(len(index.BWT))

	for i := len(pattern) - 1; i >= 0; i-- {
		currChar := searchLetters[i]
		s = uint(index.c[currChar]) + index.fastRank(s-1, currChar) + 1
		e = uint(index.c[currChar]) + index.fastRank(e, currChar)
		if e <= s {
			return []uint{}
		}
	}

	results := make([]uint, e-s+1)
	for i := uint(0); i < len(results); i++ {
		results[i] = index.sa[s+i-1]
	}

	return results
}

func (index *Index) SearchForBytesBasic(pattern []byte) []int {
=======
func NewWithWaveletTree(seq *linear.Seq) *IndexFaster {
	length := seq.Len()
	suffixArray := generateSuffixArray(seq.Seq)

	index := IndexFaster{
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
			index.c[currentLetter] = uint(i)
		}
		mem = currentLetter
	}

	index.wtree = wavelettree.New(index.BWT)

	return &index
}

func (index *Index) SearchForBytes(pattern []byte) []int {
>>>>>>> 8c03a7ded2ff0b2901294810386cdd046be11ab4
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

func (index *Index) fastRank(ranking uint, c alphabet.Letter) uint {
	// TODO Use wavelet trees to get fast counts.
	return index.wtree.Rank(uint(ranking), byte(c))
}
