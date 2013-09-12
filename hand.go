// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"sort"
	"strings"
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
// number of cards (as specified at construction time) will cause a panic.
func (hand *Hand) DealTo(card Card) {
	if hand.pos >= len(hand.cards) {
		panic("attempt to deal too many cards to hand")
	}
	hand.cards[hand.pos] = card
	hand.pos++
}

// Evaluate this hand's rank.
func (hand *Hand) Evaluate() {
	if hand.Size() < 1 {
		return
	}
	cards := hand.copyAndSortCards()
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

// What is this hand's rank?
func (hand *Hand) Rank() HandRank {
	return hand.rank
}

// Return string representation of this set of cards.
func printHand(cards []Card) string {
	rv := make([]string, len(cards))
	for i, card := range cards {
		rv[i] = card.ToString()
	}
	return strings.Join(rv, ",")
}

// ----- INTERNAL FUNCTIONS --------------------------------------------------

// Copy this hand's cards into a new array and sort them according to the
// Cards interface.
func (hand *Hand) copyAndSortCards() []Card {
	t := make([]Card, hand.Size())
	copy(t, hand.cards)
	sort.Sort(Cards(t))
	return t
}

// Does this set of cards contain a straight flush?
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

// Does this set of cards contain four of a kind?
func hasQuads(cards []Card) bool {
	return hasRunOfSameRank(cards, 4)
}

// Does this set of cards contain a full house?
func hasFullHouse(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	trips := 0
	pairs := 0
	inARow := 1
	for i := 1; i < len(cards); i++ {
		diff := byte(cards[i].Rank()) - byte(cards[i-1].Rank())
		if diff == 0 {
			inARow++
			if inARow == 3 {
				trips++
			} else if inARow == 2 {
				pairs++
			}
		} else {
			inARow = 1
		}
	}
	return trips > 1 || (trips == 1 && pairs > 1)
}

// Does this set of cards contain a flush?
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

// Does this set of cards contain a straight?
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

// Does this set of cards contain a set? (three of a kind)
func hasSet(cards []Card) bool {
	return hasRunOfSameRank(cards, 3)
}

// Does this set of cards contain two pair?
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

// Does this set of cards contain a pair?
func hasPair(cards []Card) bool {
	return hasRunOfSameRank(cards, 2)
}

// Does this set of cards have a run of runLength cards with the same rank?
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
