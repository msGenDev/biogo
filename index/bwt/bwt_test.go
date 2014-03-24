package bwt

import (
	"code.google.com/p/biogo/alphabet"
	"code.google.com/p/biogo/io/seqio/fasta"
	"code.google.com/p/biogo/seq/linear"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func demoSequence(sequence string) *linear.Seq {
	return linear.NewSeq("example DNA", []alphabet.Letter(sequence), alphabet.DNA)
}

func TestSuffixArray(t *testing.T) {
	testSeq := "TAGCTACTGATGCGTAGCTATGCTAGC"
	d := demoSequence(testSeq)

	Convey("Given a new index", t, func() {
		index := New(d)
		Convey("The BWT and suffix arrays should have the same length as the original sequence", func() {
			So(len(index.BWT), ShouldEqual, len(testSeq))
			So(len(index.sa), ShouldEqual, len(testSeq))
		})

		Convey("The suffix array should be calculated correctly", func() {
			expected := []int{5, 24, 1, 15, 9, 19, 26, 12, 3, 22, 17, 6, 8, 25, 11, 2, 21, 16, 13, 4, 23, 0, 14, 18, 7, 10, 20}
			for i := range expected {
				So(index.sa[i], ShouldEqual, expected[i])
			}
		})
	})
}

func TestIndexSearch(t *testing.T) {
	testSeq := "TAGCTACTGATGCGTAGCTATGCTAGC"
	d := demoSequence(testSeq)
	Convey("Given a new index of the string 'TAGCTACTGATGCGTAGCTATGCTAGC'", t, func() {
		index := New(d)
		Convey("When searching for 'TGA'", func() {
			results := index.SearchForBytes([]byte("TAG"))
			Convey("There should be three hits", func() {
				So(len(results), ShouldEqual, 3)
			})

			Convey("The hits should be at 0, 14, and 23.", func() {
				So(results, ShouldContain, 0)
				So(results, ShouldContain, 14)
				So(results, ShouldContain, 23)
			})
		})

		Convey("When searching for 'GTAG'", func() {
			results := index.SearchForBytes([]byte("TAGC"))
			Convey("There should be at least one hit", func() {
				So(len(results), ShouldBeGreaterThan, 0)
			})
			Convey("The hits should be at 0 and 23", func() {
				So(results, ShouldContain, 0)
				So(results, ShouldContain, 23)
			})
		})

		Convey("When searching for nonsense", func() {
			results := index.SearchForBytes([]byte("foobar"))
			Convey("There should be no results", func() {
				So(len(results), ShouldEqual, 0)
			})
		})
	})
}

func TestIndexSearchFaster(t *testing.T) {
	testSeq := "TAGCTACTGATGCGTAGCTATGCTAGC"
	d := demoSequence(testSeq)
	Convey("Given a new index of the string 'TAGCTACTGATGCGTAGCTATGCTAGC'", t, func() {
		index := NewWithWaveletTree(d)

		Convey("The index should calculate the number of 'T's in the first 10 bp", func() {
			count := index.wtree.Rank(10, 'T')
			So(count, ShouldEqual, 5)
		})
	})

	Convey("Given a new index of the string 'TAGCTACTGATGCGTAGCTATGCTAGC'", t, func() {
		index := NewWithWaveletTree(d)

		Convey("When searching for 'TGA'", func() {
			results := index.SearchForBytes([]byte("TAG"))
			Convey("There should be three hits", func() {
				So(len(results), ShouldEqual, 3)
			})

			Convey("The hits should be at 0, 14, and 23.", func() {
				So(results, ShouldContain, 0)
				So(results, ShouldContain, 14)
				So(results, ShouldContain, 23)
			})
		})

		Convey("When searching for 'GTAG'", func() {
			results := index.SearchForBytes([]byte("TAGC"))
			Convey("There should be at least one hit", func() {
				So(len(results), ShouldBeGreaterThan, 0)
			})
			Convey("The hits should be at 0 and 23", func() {
				So(results, ShouldContain, 0)
				So(results, ShouldContain, 23)
			})
		})

		Convey("When searching for nonsense", func() {
			results := index.SearchForBytes([]byte("foobar"))
			Convey("There should be no results", func() {
				So(len(results), ShouldEqual, 0)
			})
		})
	})
}

func BenchmarkIndexCreation(b *testing.B) {
	testSeq := "TAGCTACTGATGCGTAGCTATGCTAGC"
	d := demoSequence(testSeq)
	var index *Index
	for i := 0; i < b.N; i++ {
		index = New(d)
	}

	if len(index.BWT) != 27 {
		b.Fail()
	}
}

func BenchmarkFastIndexCreation(b *testing.B) {
	testSeq := "TAGCTACTGATGCGTAGCTATGCTAGC"
	d := demoSequence(testSeq)
	var index *IndexFaster
	for i := 0; i < b.N; i++ {
		index = NewWithWaveletTree(d)
	}

	if len(index.BWT) != 27 {
		b.Fail()
	}
}

func BenchmarkIndexSearch(b *testing.B) {
	f, _ := os.Open("data/subseq.fasta")
	reader := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNA))
	d, _ := reader.Read()
	s := d.(*linear.Seq)
	index := New(s)
	var results []int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = index.SearchForBytes([]byte("TAG"))
	}

	if len(results) < 100 {
		b.Errorf("Got %d results, expected at least 100\n", len(results))
	}
}

func BenchmarkFastIndexSearch(b *testing.B) {
	f, _ := os.Open("data/subseq.fasta")
	reader := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNA))
	d, _ := reader.Read()
	s := d.(*linear.Seq)
	index := NewWithWaveletTree(s)
	var results []int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = index.SearchForBytes([]byte("TAG"))
	}

	if len(results) < 100 {
		b.Errorf("Got %d results, expected at least 100\n", len(results))
	}
}