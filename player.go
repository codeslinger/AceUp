// vim:set ts=2 noet ai ft=go:
// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

const (
	MICROS_PER_US_CENT = 10000
	MICROS_PER_USD     = 100 * MICROS_PER_US_CENT
)

var idStamp = 1

type Player struct {
	id       uint64
	nick     string
	bankroll uint64
}

func NewPlayer(n string, initialBankroll uint64) *Player {
	player := new(Player)
	player.nick = n
	player.bankroll = initialBankroll
	return player
}

func (player *Player) Credit(amount uint64) {
	player.bankroll += amount
}

func (player *Player) Debit(amount uint64) {
	if amount > player.bankroll {
		player.bankroll = 0
		return
	}
	player.bankroll -= amount
}

func (player *Player) ChangeNick(newNick string) {
	player.nick = newNick
}
