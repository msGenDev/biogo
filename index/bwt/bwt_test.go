package bwt

import "testing"
import "code.google.com/p/biogo/seq/linear"
import "code.google.com/p/biogo/alphabet"

func demoSequence(sequence string) *linear.Seq {
	return linear.NewSeq("example DNA", []alphabet.Letter(sequence), alphabet.DNA)
}

func TestSuffixArray(t *testing.T) {
	d := demoSequence("TAGCTACTGATGCGTAGCTATGCTAGC")
	index := New(d)

	expected := []int{5, 24, 1, 15, 9, 19, 26, 12, 3, 22, 17, 6, 8, 25, 11, 2, 21, 16, 13, 4, 23, 0, 14, 18, 7, 10, 20}

	if len(expected) != len(index.BWT) {
		t.Errorf("BWA index failed")
	}
	for i := range expected {
		if index.sa[i] != expected[i] {
			t.Errorf("BWA suffix array construction failed")
		}
	}
}

func TestIndexSearch(t *testing.T) {
	var results []int
	d := demoSequence("TGATAGCTACTGATGCGTAGCTATGCTAGCTAGTCTAGCTGACTAGCTA")
	index := New(d)

	results = index.SearchForBytesBasic([]byte("TGA"))
	checkResultsMatch(t, results, []int{0, 10, 39})

	results = index.SearchForBytesBasic([]byte("AATGAT"))
	checkResultsMatch(t, results, []int{})

	results = index.SearchForBytesBasic([]byte("This should not match."))
	checkResultsMatch(t, results, []int{})

	results = index.SearchForBytesBasic([]byte(""))
	checkResultsMatch(t, results, index.sa)
}

func checkResultsMatch(t *testing.T, expected, results []int) {
	if len(expected) != len(results) {
		t.Errorf("Search is not correct. Expected %v, got %v", expected, results)
	}

	for i := range expected {
		absent := true
		for j := range results {
			if expected[i] == results[j] {
				absent = false
			}
		}
		if absent {
			t.Errorf("Search is not correct. Expected %v, got %v", expected, results)
		}
	}
}
