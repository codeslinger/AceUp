// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

const (
	MICROS_PER_US_CENT = 10000
	MICROS_PER_USD     = 100 * MICROS_PER_US_CENT
)

var idStamp = 1

// ----- Player Public API ---------------------------------------------------

type Player interface {
	Bankroll() uint64
	Credit(uint64)
	Debit(uint64)
	Nick() string
	ChangeNick(string)
}

func NewPlayer(n string, initialBankroll uint64) Player {
	return &player{
		nick:     n,
		bankroll: initialBankroll,
	}
}

func (p *player) Bankroll() uint64 {
	return p.bankroll
}

func (p *player) Credit(amount uint64) {
	p.bankroll += amount
}

func (p *player) Debit(amount uint64) {
	if amount > p.bankroll {
		p.bankroll = 0
		return
	}
	p.bankroll -= amount
}

func (p *player) Nick() string {
	return p.nick
}

func (p *player) ChangeNick(newNick string) {
	p.nick = newNick
}

// ----- Player Internal API -------------------------------------------------

type player struct {
	id       uint64
	nick     string
	bankroll uint64
}
