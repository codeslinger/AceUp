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

type Rank  byte
type Suit  byte
type Card  byte
type Cards []Card

var NoCard   = newCard(NoSuit, NoRank)
var ranks    = []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
var suits    = []Suit{Club, Diamond, Heart, Spade}
var rankStr  = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
var suitStr  = []string{"c", "d", "h", "s"}

// ----- SUIT PUBLIC API ----------------------------------------------------

// Is this suit valid?
func (suit Suit) Valid() bool {
  return !(suit < Club || suit > Spade)
}

// ----- RANK PUBLIC API ----------------------------------------------------

// Is this rank valid?
func (rank Rank) Valid() bool {
  return !(rank < Two || rank > Ace)
}

// ----- CARDS PUBLIC API ---------------------------------------------------

// All these implement the sort.Interface interface so we can sort an array
// of Card values

func (cards Cards) Len() int {
  return len(cards)
}

func (cards Cards) Swap(i, j int) {
  cards[i], cards[j] = cards[j], cards[i]
}

func (cards Cards) Less(i, j int) bool {
  return cards[i].Compare(cards[j]) < 0
}

// ----- CARD PUBLIC API ----------------------------------------------------

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

// Compare two cards by rank only, returning:
//   >0 if this card rank is higher
//   <0 if the other card rank is higher, or
//    0 if they represent the same card rank.
func (card Card) Compare(other Card) int {
  return int(card.Rank() - other.Rank())
}

// Compare two cards by rank first and then suit, returning:
//   >0 if this card is higher
//   <0 if the other card is higher, or
//    0 if they represent the same card.
func (card Card) FullCompare(other Card) int {
  cmp := card.Compare(other)
  if cmp != 0 {
    return cmp
  }
  return int(card.Suit() - other.Suit())
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
  card = Card(((byte(rank) << 4) & 0xF0) | (byte(suit) & 0x0F))
  return
}

// NB: we put the rank first in the byte and then the suit so the cards will
// sort properly when converted to integral values

func (card Card) rank() Rank {
  return ranks[((byte(card) >> 4) & 0x0F) - 1]
}

func (card Card) suit() Suit {
  return suits[(byte(card) & 0x0F) - 1]
}

