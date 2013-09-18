// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package poker

import (
	"fmt"
	"strings"
)

const CardsPerDeck = 52

// Hand rankings
const (
	StraightFlush = 1
	FourOfAKind   = 2
	FullHouse     = 3
	Flush         = 4
	Straight      = 5
	ThreeOfAKind  = 6
	TwoPair       = 7
	OnePair       = 8
	HighCard      = 9
)

// Card suits
const (
	Club    = 0x8000
	Diamond = 0x4000
	Heart   = 0x2000
	Spade   = 0x1000
)

// Card ranks
const (
	Deuce = 0
	Trey  = 1
	Four  = 2
	Five  = 3
	Six   = 4
	Seven = 5
	Eight = 6
	Nine  = 7
	Ten   = 8
	Jack  = 9
	Queen = 10
	King  = 11
	Ace   = 12
)

var rankStr = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}

// ----- CARD API ------------------------------------------------------------

type Card uint32

// Return the Card representing the given rank and suit.
func NewCard(rank int, suit int) Card {
	r := uint32(rank)
	s := uint32(suit)
	return Card(primes[r] | (r << 8) | s | (1 << (16 + r)))
}

// Extract and report the rank of this card.
func (card Card) Rank() int {
	return int((card >> 8) & 0xF)
}

// Extract and report the suit of this card.
func (card Card) Suit() int {
	if card&Club != 0 {
		return Club
	}
	if card&Diamond != 0 {
		return Diamond
	}
	if card&Heart != 0 {
		return Heart
	}
	return Spade
}

// Return the string representation of this card. (e.g. "Td", "As", "9c", etc)
func (card Card) String() string {
	var suit string

	switch card.Suit() {
	case Club:
		suit = "c"
	case Diamond:
		suit = "d"
	case Heart:
		suit = "h"
	case Spade:
		suit = "s"
	}
	return fmt.Sprintf("%s%s", rankStr[card.Rank()], suit)
}

// ----- DECK API ------------------------------------------------------------

var EmptyDeck = fmt.Errorf("deck is empty")

type Deck interface {
	Shuffle()
	Deal() Card
	Empty() bool
	Remaining() int
}

type PokerDeck struct {
	cards []Card
	pos   int
}

// Create a new deck of cards. This deck will *NOT* be shuffled.
func NewPokerDeck() *PokerDeck {
	deck := &PokerDeck{cards: make([]Card, CardsPerDeck)}
	n := 0
	for suit := Club; suit >= Spade; suit >>= 1 {
		for rank := Deuce; rank <= Ace; rank++ {
			deck.cards[n] = NewCard(rank, suit)
			n++
		}
	}
	return deck
}

// Shuffle this deck of cards. We use a Fisher-Yates shuffle:
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle
func (deck *PokerDeck) Shuffle() {
	for i := len(deck.cards) - 1; i > 0; i-- {
		if j := randInt(i + 1); j >= 0 && i != j {
			deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
		}
	}
	deck.pos = 0
}

// Deal one card from the deck. Returns error on empty deck.
func (deck *PokerDeck) Deal() Card {
	if deck.Empty() {
		panic(EmptyDeck)
	}
	card := deck.cards[deck.pos]
	deck.pos++
	return card
}

// Is this deck empty?
func (deck *PokerDeck) Empty() bool {
	return deck.Remaining() == 0
}

// How many cards are remaining in this deck?
func (deck *PokerDeck) Remaining() int {
	return (len(deck.cards) - deck.pos) - 1
}

// ----- PUBLIC HAND EVALUATION API ------------------------------------------

// Determine the given hand's ranking. (i.e. StraightFlush, ThreeOfAKind, etc)
func EvaluateForHigh(hand []Card) int {
	return handRank(eval7CardHand(hand))
}

// Return string representation of hand of cards.
func PrintHand(hand []Card) string {
	rv := make([]string, len(hand))
	for i, card := range hand {
		rv[i] = card.String()
	}
	return "(" + strings.Join(rv, ",") + ")"
}

// ----- HAND EVALUATION FUNCTIONS -------------------------------------------

// Determine a hand's rank given an equivalence value.
func handRank(val uint16) int {
	if val > 6185 {
		return HighCard
	}
	if val > 3325 {
		return OnePair
	}
	if val > 2467 {
		return TwoPair
	}
	if val > 1609 {
		return ThreeOfAKind
	}
	if val > 1599 {
		return Straight
	}
	if val > 322 {
		return Flush
	}
	if val > 166 {
		return FullHouse
	}
	if val > 10 {
		return FourOfAKind
	}
	return StraightFlush
}

// Generate the equivalence value for a 7-card hand. This is unoptimized, as
// it will evaluate all possible 5-card hands in a 7-card hand (7-choose-5,
// or 21) and return the best equivalence value found.
func eval7CardHand(hand []Card) uint16 {
	var best uint16 = 0xFFFF

	subhand := []Card{0, 0, 0, 0, 0}
	for i := 0; i < 21; i++ {
		for j := 0; j < 5; j++ {
			subhand[j] = hand[perm7[i][j]]
		}
		q := eval5CardHand(subhand)
		if q < best {
			best = q
		}
	}
	return best
}

// Generate the equivalence value for a 5-card hand.
func eval5CardHand(hand []Card) uint16 {
	return eval5CardHandFast(uint32(hand[0]), uint32(hand[1]), uint32(hand[2]), uint32(hand[3]), uint32(hand[4]))
}

// Generate the equivalence value for a sampling of 5 cards.
func eval5CardHandFast(c1, c2, c3, c4, c5 uint32) uint16 {
	q := (c1 | c2 | c3 | c4 | c5) >> 16
	// check for flush/straight flush
	if (c1 & c2 & c3 & c4 & c5 & 0xF000) != 0 {
		return flushes[q]
	}
	// check for straight/high card
	s := unique5[q]
	if s != 0 {
		return s
	}
	i := findFast((c1 & 0xFF) * (c2 & 0xFF) * (c3 & 0xFF) * (c4 & 0xFF) * (c5 & 0xFF))
	return hash_values[i]
}

func findFast(u uint32) uint32 {
	u += 0xE91AAA35
	u ^= u >> 16
	u += u << 8
	u ^= u >> 4
	b := (u >> 8) & 0x1FF
	a := (u + (u << 2)) >> 19
	return a ^ hash_adjust[b]
}
