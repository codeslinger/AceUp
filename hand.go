// vim:set ts=2 noet ai ft=go:
// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"sort"
)

// hand ranks in ascending order (i.e. StraightFlush beats all)
const (
	NoHandRank    = 0
	HighCard      = 1
	OnePair       = 2
	TwoPair       = 3
	ThreeOfAKind  = 4
	Straight      = 5
	Flush         = 6
	FullHouse     = 7
	FourOfAKind   = 8
	StraightFlush = 9
)

const HAND_SIZE = 5

type HandRank byte

type Hand struct {
	cards []Card
	pos   int
	rank  HandRank
}

var handRanks = []HandRank{
	HighCard,
	OnePair,
	TwoPair,
	ThreeOfAKind,
	Straight,
	Flush,
	FullHouse,
	FourOfAKind,
	StraightFlush,
}

// ----- HAND PUBLIC API -----------------------------------------------------

// Construct a new Hand instance.
func NewHand(maxCards int) *Hand {
	hand := new(Hand)
	hand.rank = NoHandRank
	hand.cards = make([]Card, maxCards)
	hand.pos = 0
	return hand
}

// Is this hand higher than, lower than or equal to another hand?
// Returns:
//   >0 if this hand is higher
//   <0 if the other hand is higher
//    0 if they are equal.
func (hand *Hand) Compare(other *Hand) int {
	// TODO: implement this
	return 0
}

// How many cards are in this hand?
func (hand *Hand) Size() int {
	return hand.pos
}

// Deal this hand a card. Attempts to deal this hand more than its maximum
// number of cards (as specified at construction time) will be silently
// ignored.
func (hand *Hand) DealTo(card Card) {
	if hand.pos >= len(hand.cards) {
		return
	}
	hand.cards[hand.pos] = card
	hand.pos++
}

// Evaluate this hand's rank.
func (hand *Hand) Evaluate() {
	var cards []Card

	if cards = hand.copyAndSortCards(); cards == nil {
		return
	}
	if hasStraightFlush(cards) {
		hand.rank = StraightFlush
	} else if hasQuads(cards) {
		hand.rank = FourOfAKind
	} else if hasFullHouse(cards) {
		hand.rank = FullHouse
	} else if hasFlush(cards) {
		hand.rank = Flush
	} else if hasStraight(cards) {
		hand.rank = Straight
	} else if hasSet(cards) {
		hand.rank = ThreeOfAKind
	} else if hasTwoPair(cards) {
		hand.rank = TwoPair
	} else if hasPair(cards) {
		hand.rank = OnePair
	} else {
		hand.rank = HighCard
	}
}

// ----- INTERNAL FUNCTIONS --------------------------------------------------

func (hand *Hand) copyAndSortCards() []Card {
	if hand.Size() < 1 {
		return nil
	}
	t := make([]Card, hand.Size())
	copy(hand.cards, t)
	sort.Sort(Cards(t))
	return t
}

func hasStraightFlush(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	inARow := 0
	for i := range cards {
		if i > 0 {
			rankDiff := byte(cards[i].Rank()) - byte(cards[i-1].Rank())
			suitDiff := byte(cards[i].Suit()) - byte(cards[i-1].Suit())
			if rankDiff == 1 && suitDiff == 0 {
				inARow++
				if inARow == 4 && cards[i].Rank() == Two && cards[0].Rank() == Ace && cards[0].Suit() == cards[i].Suit() {
					// found a Steal Wheel (A-2-3-4-5 of same suit)
					inARow++
				}
			} else if rankDiff != 0 {
				inARow = 1
			}
		}
		if inARow == HAND_SIZE {
			return true
		}
	}
	return false
}

func hasQuads(cards []Card) bool {
	return hasRunOfSameRank(cards, 4)
}

func hasFullHouse(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	trips := 0
	pairs := 0
	inARow := 1
	for i := range cards {
		diff := byte(cards[i].Rank()) - byte(cards[i-1].Rank())
		if diff == 0 {
			inARow++
			if inARow == 3 {
				trips++
			} else if inARow == 2 {
				pairs++
			}
		} else {
			inARow = 0
		}
	}
	return trips > 1 || (trips == 1 && pairs > 1)
}

func hasFlush(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	suitsFound := []byte{0, 0, 0, 0}
	for i := range cards {
		idx := int(cards[i].Suit()) - int(Club)
		suitsFound[idx]++
		if suitsFound[idx] == HAND_SIZE {
			return true
		}
	}
	return false
}

func hasStraight(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	inARow := 1
	for i := range cards {
		if i == 0 {
			continue
		}
		diff := byte(cards[i].Rank()) - byte(cards[i-1].Rank())
		if diff == 1 {
			inARow++
			if inARow == 4 && cards[i].Rank() == Two && cards[0].Rank() == Ace {
				// found a Wheel (A-2-3-4-5)
				inARow++
			}
		} else if diff != 0 {
			inARow = 1
		}
		if inARow == HAND_SIZE {
			return true
		}
	}
	return false
}

func hasSet(cards []Card) bool {
	return hasRunOfSameRank(cards, 3)
}

func hasTwoPair(cards []Card) bool {
	if len(cards) < 4 {
		return false
	}
	pairs := 0
	inARow := 1
	for i := range cards {
		if i == 0 {
			continue
		}
		diff := byte(cards[i].Rank()) - byte(cards[i-1].Rank())
		if diff == 0 {
			inARow++
			if inARow == 2 {
				pairs++
				inARow = 1
			}
		} else {
			inARow = 0
		}
	}
	return pairs >= 2
}

func hasPair(cards []Card) bool {
	return hasRunOfSameRank(cards, 2)
}

func hasRunOfSameRank(cards []Card, runLength int) bool {
	if len(cards) < runLength {
		return false
	}
	inARow := 1
	for i := range cards {
		if i == 0 {
			continue
		}
		diff := byte(cards[i].Rank()) - byte(cards[i-1].Rank())
		if diff == 0 {
			inARow++
			if inARow == runLength {
				return true
			}
		} else {
			inARow = 1
		}
	}
	return false
}
