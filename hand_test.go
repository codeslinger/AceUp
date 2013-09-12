// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"testing"
)

func Test_hand_size_calculated_correctly(t *testing.T) {
	hand := NewHand(7)
	if hand.Size() != 0 {
		t.Fatalf("hand size should be 0 but was %d", hand.Size())
	}
	for i, rank := range []Rank{Ace, Two, Three, Four, Five, Six, Seven} {
		hand.DealTo(MakeCard(t, rank, Spade))
		if hand.Size() != i+1 {
			t.Fatalf("hand size should be %d but was %d", i+1, hand.Size())
		}
	}
}

func Test_can_detect_straight_flush(t *testing.T) {
	hand := makeHand(t, []string{"Ts", "9s", "7s", "Js", "Ad", "3c", "8s"})
	hand.Evaluate()
	if hand.Rank() != StraightFlush {
		t.Fatalf("expected %d but was %d", StraightFlush, hand.Rank())
	}
}

func Test_can_detect_four_of_a_kind(t *testing.T) {
	hand := makeHand(t, []string{"Ts", "Tc", "8h", "7s", "Td", "Kd", "Th"})
	hand.Evaluate()
	if hand.Rank() != FourOfAKind {
		t.Fatalf("expected %d but was %d", FourOfAKind, hand.Rank())
	}
}

func Test_can_detect_full_house(t *testing.T) {
	hand := makeHand(t, []string{"Ts", "Tc", "8h", "7s", "Td", "Kd", "Kh"})
	hand.Evaluate()
	if hand.Rank() != FullHouse {
		t.Fatalf("expected %d but was %d", FullHouse, hand.Rank())
	}
}

func Test_can_detect_full_house_in_two_sets(t *testing.T) {
	hand := makeHand(t, []string{"Ts", "Tc", "8h", "Ks", "Td", "Kd", "Kh"})
	hand.Evaluate()
	if hand.Rank() != FullHouse {
		t.Fatalf("expected %d but was %d", FullHouse, hand.Rank())
	}
}

func Test_can_detect_flush(t *testing.T) {
	hand := makeHand(t, []string{"Ts", "9s", "8h", "Ks", "7s", "2s", "Kh"})
	hand.Evaluate()
	if hand.Rank() != Flush {
		t.Fatalf("expected %d but was %d", Flush, hand.Rank())
	}
}

func makeHand(t *testing.T, cards []string) *Hand {
	hand := NewHand(len(cards))
	for _, str := range cards {
		card, err := CardFromString(str)
		if err != nil {
			t.Fatalf("bad card specification: '%s'", str)
		}
		hand.DealTo(card)
	}
	return hand
}
