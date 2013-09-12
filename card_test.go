// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"testing"
)

func Test_creates_card_with_correct_rank(t *testing.T) {
	card := MakeCard(t, Queen, Diamond)
	if card.Rank() != Queen {
		t.Fatalf("Rank does not match: expected=%d actual=%d (card=%d)", Queen, card.Rank(), card)
	}
}

func Test_creates_card_with_correct_suit(t *testing.T) {
	card := MakeCard(t, Queen, Diamond)
	if card.Suit() != Diamond {
		t.Fatalf("Suit does not match: expected=%d actual=%d (card=%d)", Diamond, card.Suit(), card)
	}
}

func Test_can_create_smallest_card(t *testing.T) {
	card := MakeCard(t, Two, Club)
	if card.Rank() != Two || card.Suit() != Club {
		t.Fatalf("Card does not match: expected=0 actual=%d", card)
	}
}

func Test_can_create_biggest_card(t *testing.T) {
	card := MakeCard(t, Ace, Spade)
	if card.Rank() != Ace || card.Suit() != Spade {
		t.Fatalf("Card does not match: expected=51 actual=%d", card)
	}
}

func Test_higher_rank_yields_negative_comparison_value(t *testing.T) {
	a := MakeCard(t, Three, Club)
	b := MakeCard(t, Four, Club)
	rv := a.Compare(b)
	if rv >= 0 {
		t.Fatalf("Expected negative comparison value; got %d", rv)
	}
}

func Test_lower_rank_yields_positive_comparison_value(t *testing.T) {
	a := MakeCard(t, Four, Club)
	b := MakeCard(t, Three, Club)
	rv := a.Compare(b)
	if rv <= 0 {
		t.Fatalf("Expected positive comparison value; got %d", rv)
	}
}

func Test_of_same_rank_cards_yields_zero_comparison_value(t *testing.T) {
	a := MakeCard(t, Three, Diamond)
	b := MakeCard(t, Three, Heart)
	rv := a.Compare(b)
	if rv != 0 {
		t.Fatalf("Expected zero comparison value; got %d", rv)
	}
}

func Test_string_representations_of_cards_are_correct(t *testing.T) {
	ranks := []Rank{Three, Seven, Ten, King}
	suits := []Suit{Club, Diamond, Heart, Spade}
	strings := []string{"3c", "7d", "Th", "Ks"}
	for i := range ranks {
		card := MakeCard(t, ranks[i], suits[i])
		if card.ToString() != strings[i] {
			t.Fatalf("Incorrect string repr: expected=%s actual=%s", strings[i], card)
		}
	}
}

func MakeCard(t *testing.T, rank Rank, suit Suit) Card {
	card, err := NewCard(rank, suit)
	if err != nil {
		t.Fatalf("Error creating test card: %v", err)
	}
	return card
}
