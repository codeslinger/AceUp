// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"fmt"
)

const MAX_CARD = 52

// card suits
const (
	Club    = 1
	Diamond = 2
	Heart   = 3
	Spade   = 4
)

// card ranks
const (
	Two   = 1
	Three = 2
	Four  = 3
	Five  = 4
	Six   = 5
	Seven = 6
	Eight = 7
	Nine  = 8
	Ten   = 9
	Jack  = 10
	Queen = 11
	King  = 12
	Ace   = 13
)

type Rank byte
type Suit byte
type Card byte
type Cards []Card

var (
	NoCard  = Card(byte(255))
	rankStr = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	suitStr = []string{"c", "d", "h", "s"}
)

// ----- CARDS PUBLIC API ----------------------------------------------------

// Implement sort.Interface so we can sort Cards
func (c Cards) Len() int           { return len(c) }
func (c Cards) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Cards) Less(i, j int) bool { return c[i].Compare(c[j]) < 0 }

// ----- CARD PUBLIC API -----------------------------------------------------

// Create a new Card instance.
func NewCard(rank Rank, suit Suit) (Card, error) {
	if rank < Two || rank > Ace {
		return 0, fmt.Errorf("invalid rank specified")
	}
	if suit < Club || suit > Spade {
		return 0, fmt.Errorf("invalid suit specified")
	}
	return Card((byte(suit)-1)*13 + byte(rank) - 1), nil
}

// Get the rank of this card.
func (card Card) Rank() Rank {
	return Rank((byte(card) % 13) + 1)
}

// Get the suit of this card.
func (card Card) Suit() Suit {
	return Suit((byte(card) / 13) + 1)
}

// Compare two cards by rank only, returning:
//   >0 if this card rank is higher
//   <0 if the other card rank is higher, or
//    0 if they represent the same card rank.
func (card Card) Compare(other Card) int {
	return int(card.Rank()) - int(other.Rank())
}

// Return string representation of this card.
func (card Card) ToString() string {
	return fmt.Sprintf("%s%s", rankStr[byte(card)%13], suitStr[byte(card)/13])
}
