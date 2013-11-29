// Copyright ©2011-2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package align

import (
	"code.google.com/p/biogo/alphabet"
	"code.google.com/p/biogo/seq/linear"

	"fmt"
)

func ExampleFitted_Align() {
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte("AGACTAGATT"))}
	fsa.Alpha = alphabet.DNA
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte("GACAGACGA"))}
	fsb.Alpha = alphabet.DNA

	//		   Query letter
	//  	 A	 C	 G	 T	 -
	// A	10	-3	-1	-4	-5
	// C	-3	 9	-5	 0	-5
	// G	-1	-5	 7	-3	-5
	// T	-4	 0	-3	 8	-5
	// -	-5	-5	-5	-5	 0
	fitted := Fitted{
		{10, -3, -1, -4, -5},
		{-3, 9, -5, 0, -5},
		{-1, -5, 7, -3, -5},
		{-4, 0, -3, 8, -5},
		{-5, -5, -5, -5, 0},
	}

	aln, err := fitted.Align(fsa, fsb)
	if err == nil {
		fmt.Printf("%s\n", aln)
		fa := Format(fsa, fsb, aln, '-')
		fmt.Printf("%s\n%s\n", fa[0], fa[1])
	}
	// Output:
	//[[1,4)/[0,3)=26 [4,5)/-=-5 [5,10)/[3,8)=24]
	// GACTAGATT
	// GAC-AGACG
}

func ExampleFittedAffine_Align() {
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte("ATAGGAA"))}
	fsa.Alpha = alphabet.DNA
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte("ATTGGCAATGA"))}
	fsb.Alpha = alphabet.DNA

	//		   Query letter
	//  	 A	 C	 G	 T	 -
	// A	 1	-1	-1	-1	-1
	// C	-1	 1	-1	 1	-1
	// G	-1	-1	 1	-1	-1
	// T	-1	 1	-1	 1	-1
	// -	-1	-1	-1	-1	 0
	//
	// Gap open: -5
	fitted := FittedAffine{
		Matrix: Linear{
			{1, -1, -1, -1, -1},
			{-1, 1, -1, -1, -1},
			{-1, -1, 1, -1, -1},
			{-1, -1, -1, 1, -1},
			{-1, -1, -1, -1, 0},
		},
		GapOpen: -5,
	}

	aln, err := fitted.Align(fsa, fsb)
	if err == nil {
		fmt.Printf("%s\n", aln)
		fa := Format(fsa, fsb, aln, '-')
		fmt.Printf("%s\n%s\n", fa[0], fa[1])
	}
	// Output:
	// [[0,7)/[0,7)=3]
	// ATAGGAA
	// ATTGGCA
}
