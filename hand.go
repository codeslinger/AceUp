// vim:set ts=2 sw=2 et ai ft=go:
// Copyright (c) 2012 Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

// hand ranks
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

type HandRank byte

type Hand struct {
  cards []Card
}

var handRanks = []HandRank{HighCard, OnePair, TwoPair, ThreeOfAKind, Straight, Flush, FullHouse, FourOfAKind, StraightFlush}

func NewHand() *Hand {
  hand := new(Hand)
  return hand
}

func (hand *Hand) Rank() HandRank {
  return NoHandRank
}

func (hand *Hand) isStraightFlush() bool {
  return false
}

