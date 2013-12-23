package bwt

import (
	"code.google.com/p/biogo/alphabet"
	"code.google.com/p/biogo/seq/linear"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func demoSequence(sequence string) *linear.Seq {
	return linear.NewSeq("example DNA", []alphabet.Letter(sequence), alphabet.DNA)
}

func TestIndexSearch(t *testing.T) {
	testSeq := "TAGCTACTGATGCGTAGCTATGCTAGC"
	d := demoSequence(testSeq)
	Convey("Given a new index of the string 'TAGCTACTGATGCGTAGCTATGCTAGC'", t, func() {
		index := New(d)
		Convey("When searching for 'TGA'", func() {
			results := index.SearchForBytesBasic([]byte("TAG"))
			Convey("There should be three hits", func() {
				So(len(results), ShouldEqual, 3)
			})

			Convey("The hits should be at 0, 14, and 23.", func() {
				So(results, ShouldContain, 0)
				So(results, ShouldContain, 14)
				So(results, ShouldContain, 23)
			})
		})

		Convey("When searching for a pattern that wraps around the string", func() {
			results := index.SearchForBytesBasic([]byte("AGCTAGCTACT"))
			Convey("There should be no results returned", func() {
				So(len(results), ShouldBeZeroValue)
			})
		})

		Convey("When searching for 'GTAG'", func() {
			results := index.SearchForBytesBasic([]byte("TAGC"))
			Convey("There should be at least one hit", func() {
				So(len(results), ShouldBeGreaterThan, 0)
			})
			Convey("The hits should be at 0 and 23", func() {
				So(results, ShouldContain, 0)
				So(results, ShouldContain, 23)
			})
		})

		Convey("When searching for nonsense", func() {
			results := index.SearchForBytesBasic([]byte("foobar"))
			Convey("There should be no results", func() {
				So(len(results), ShouldBeZeroValue)
			})
		})
	})
}
