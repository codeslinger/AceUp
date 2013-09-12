// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"crypto/rand"
	"math/big"
)

type Deck struct {
	pos   int
	cards [MAX_CARD]Card
}

// ----- DECK PUBLIC API -----------------------------------------------------

// Create a new deck of cards.
func NewDeck() *Deck {
	deck := new(Deck)
	i := 0
	for suit := Club; suit <= Spade; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			deck.cards[i], _ = NewCard(Rank(rank), Suit(suit))
			i++
		}
	}
	return deck
}

// Shuffle this deck of cards. It uses a Fisher-Yates shuffle:
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle
func (deck *Deck) Shuffle() {
	for i := len(deck.cards) - 1; i > 0; i-- {
		if j := randInt(i + 1); j >= 0 && i != j {
			deck.swap(i, j)
		}
	}
	deck.pos = 0
}

// Number of cards remaining that can be dealt in this Deck.
func (deck *Deck) CardsRemaining() int {
	return (len(deck.cards) - deck.pos) - 1
}

// Is this deck empty?
func (deck *Deck) Empty() bool {
	return deck.CardsRemaining() <= 0
}

// Deal the top card from this deck.
func (deck *Deck) Deal() (card Card) {
	if deck.Empty() {
		return NoCard
	}
	card = deck.cards[deck.pos]
	deck.pos++
	return card
}

// Burn a card, essentially discarding it.
func (deck *Deck) Burn() {
	if !deck.Empty() {
		deck.Deal()
	}
}

// ----- INTERNAL FUNCTIONS --------------------------------------------------

func (deck *Deck) swap(i int, j int) {
	deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
}

// randInt generates a pseudo-random number in the range [0,max). It returns
// the generated number on success or -1 if a number could not be generated or
// max was less than or equal to 0.
func randInt(max int) int {
	if max <= 0 {
		return -1
	}
	m := big.NewInt(int64(max))
	r, e := rand.Int(rand.Reader, m)
	if e != nil {
		return -1
	}
	return int(r.Int64() % int64(max))
}
