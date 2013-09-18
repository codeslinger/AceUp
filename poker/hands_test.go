// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package poker

import (
	"testing"
)

func Test_can_detect_royal_flush(t *testing.T) {
	// ace-high straight flush in diamonds
	hand := makeHand([]string{"Td", "Kd", "7s", "Jd", "Ad", "3c", "Qd"})
	rank := EvaluateForHigh(hand)
	if rank != StraightFlush {
		t.Fatalf("expected %d but was %d", StraightFlush, rank)
	}
}

func Test_can_detect_straight_flush(t *testing.T) {
	// ten-high straight flush in spades
	hand := makeHand([]string{"Ts", "9s", "7s", "Js", "Ad", "3c", "8s"})
	rank := EvaluateForHigh(hand)
	if rank != StraightFlush {
		t.Fatalf("expected %d but was %d", StraightFlush, rank)
	}
}

func Test_can_detect_steal_wheel(t *testing.T) {
	// steal wheel (A-2-3-4-5 of same suit) in clubs
	// NB: we put an off-suit ace first here to try and trick the hand evaluator
	hand := makeHand([]string{"Ad", "Ac", "2c", "4c", "Kd", "3c", "5c"})
	rank := EvaluateForHigh(hand)
	if rank != StraightFlush {
		t.Fatalf("expected %d but was %d", StraightFlush, rank)
	}
	// same hand with in-suit ace first to ensure we can catch it both ways
	hand = makeHand([]string{"Ac", "Ad", "2c", "4c", "Kd", "3c", "5c"})
	rank = EvaluateForHigh(hand)
	if rank != StraightFlush {
		t.Fatalf("expected %d but was %d", StraightFlush, rank)
	}
}

func Test_can_detect_four_of_a_kind(t *testing.T) {
	// four tens
	hand := makeHand([]string{"Ts", "Tc", "8h", "7s", "Td", "Kd", "Th"})
	rank := EvaluateForHigh(hand)
	if rank != FourOfAKind {
		t.Fatalf("expected %d but was %d", FourOfAKind, rank)
	}
}

func Test_can_detect_full_house(t *testing.T) {
	// tens full of kings
	hand := makeHand([]string{"Ts", "Tc", "8h", "7s", "Td", "Kd", "Kh"})
	rank := EvaluateForHigh(hand)
	if rank != FullHouse {
		t.Fatalf("expected %d but was %d", FullHouse, rank)
	}
}

func Test_can_detect_full_house_in_two_sets(t *testing.T) {
	// kings full of tens
	hand := makeHand([]string{"Ts", "Tc", "8h", "Ks", "Td", "Kd", "Kh"})
	rank := EvaluateForHigh(hand)
	if rank != FullHouse {
		t.Fatalf("expected %d but was %d", FullHouse, rank)
	}
}

func Test_can_detect_flush(t *testing.T) {
	// king-high flush in spades
	hand := makeHand([]string{"Ts", "9s", "8h", "Ks", "7s", "2s", "Kh"})
	rank := EvaluateForHigh(hand)
	if rank != Flush {
		t.Fatalf("expected %d but was %d", Flush, rank)
	}
}

func Test_can_detect_straight(t *testing.T) {
	// jack-high straight
	hand := makeHand([]string{"Ts", "9s", "8h", "Ks", "7s", "3d", "Jh"})
	rank := EvaluateForHigh(hand)
	if rank != Straight {
		t.Fatalf("expected %d but was %d", Straight, rank)
	}
}

func Test_can_detect_set(t *testing.T) {
	// set of nines
	hand := makeHand([]string{"Ts", "9s", "8h", "9d", "9c", "3d", "Jh"})
	rank := EvaluateForHigh(hand)
	if rank != ThreeOfAKind {
		t.Fatalf("expected %d but was %d", ThreeOfAKind, rank)
	}
}

func Test_can_detect_two_pair(t *testing.T) {
	// two-pair, tens and sevens
	hand := makeHand([]string{"Ts", "Td", "8h", "Ks", "7s", "3d", "7h"})
	rank := EvaluateForHigh(hand)
	if rank != TwoPair {
		t.Fatalf("expected %d but was %d", TwoPair, rank)
	}
}

func Test_can_detect_one_pair(t *testing.T) {
	// pair of nines
	hand := makeHand([]string{"Ts", "9s", "8h", "9c", "4c", "3d", "Jh"})
	rank := EvaluateForHigh(hand)
	if rank != OnePair {
		t.Fatalf("expected %d but was %d", OnePair, rank)
	}
}

func Test_can_detect_high_card(t *testing.T) {
	// king high
	hand := makeHand([]string{"Ts", "9s", "8h", "Ks", "4s", "3d", "Jh"})
	rank := EvaluateForHigh(hand)
	if rank != HighCard {
		t.Fatalf("expected %d but was %d", HighCard, rank)
	}
}

func makeHand(cards []string) []Card {
	hand := make([]Card, len(cards))
	for i, str := range cards {
		hand[i] = cardFromString(str)
	}
	return hand
}

func cardFromString(s string) Card {
	var rank, suit int

	switch s[0] {
	case '2':
		rank = Deuce
	case '3':
		rank = Trey
	case '4':
		rank = Four
	case '5':
		rank = Five
	case '6':
		rank = Six
	case '7':
		rank = Seven
	case '8':
		rank = Eight
	case '9':
		rank = Nine
	case 'T':
		rank = Ten
	case 'J':
		rank = Jack
	case 'Q':
		rank = Queen
	case 'K':
		rank = King
	case 'A':
		rank = Ace
	}
	switch s[1] {
	case 'c':
		suit = Club
	case 'd':
		suit = Diamond
	case 'h':
		suit = Heart
	case 's':
		suit = Spade
	}
	return NewCard(rank, suit)
}
