// vim:set ts=2 sw=2 et ai ft=go:
// Copyright (c) 2012 Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
  "fmt"
)

// card suits
const (
  NoSuit  = 0
  Club    = 1
  Diamond = 2
  Heart   = 3
  Spade   = 4
)

// card ranks
const (
  NoRank = 0
  Two    = 1
  Three  = 2
  Four   = 3
  Five   = 4
  Six    = 5
  Seven  = 6
  Eight  = 7
  Nine   = 8
  Ten    = 9
  Jack   = 10
  Queen  = 11
  King   = 12
  Ace    = 13
)

type Rank byte
type Suit byte
type Card byte

var NoCard   = newCard(NoSuit, NoRank)
var ranks    = []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
var suits    = []Suit{Club, Diamond, Heart, Spade}
var rankStr  = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
var suitStr  = []string{"c", "d", "h", "s"}

// ----- PUBLIC API ---------------------------------------------------------

// Is this suit valid?
func (suit Suit) Valid() bool {
  return !(suit < Club || suit > Spade)
}

// Is this rank valid?
func (rank Rank) Valid() bool {
  return !(rank < Two || rank > Ace)
}

// Is this card valid?
func (card Card) Valid() bool {
  return (card.rank().Valid() && card.suit().Valid())
}

// Get the rank of this card.
func (card Card) Rank() Rank {
  if card.Valid() {
    return card.rank()
  }
  return NoRank
}

// Get the suit of this card.
func (card Card) Suit() Suit {
  if card.Valid() {
    return card.suit()
  }
  return NoSuit
}

// Compare two cards, returning:
//   1 if this card is higher
//  -1 if the other card is higher, or
//   0 if they represent the same card.
func (card Card) Compare(other Card) int {
  if card.Rank() > other.Rank() {
    return 1
  } else if card.Rank() < other.Rank() {
    return -1
  } else {
    if card.Suit() > other.Suit() {
      return 1
    } else if card.Suit() < other.Suit() {
      return -1
    }
  }
  return 0
}

// Return string representation of this card.
func (card Card) ToString() string {
  return fmt.Sprintf("%s%s",
                     rankStr[card.Rank() - 1],
                     suitStr[card.Suit() - 1])
}

// ----- INTERNAL FUNCTIONS -------------------------------------------------

// NOTE: the constructor for a card is private; external code should not be
// creating cards, just using them as part of Decks and Hands
func newCard(suit Suit, rank Rank) (card Card) {
  card = Card(((byte(suit) << 4) & 0xF0) | (byte(rank) & 0x0F))
  return
}

func (card Card) rank() Rank {
  return ranks[(byte(card) & 0x0F) - 1]
}

func (card Card) suit() Suit {
  return suits[((byte(card) >> 4) & 0x0F) - 1]
}

