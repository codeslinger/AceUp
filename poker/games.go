package poker

type GameType int
const (
	// community card games
	Holdem = iota
	Omaha
	OmahaHL
	Omaha5
	Omaha5HL
	Courchevel
	CourchevelHL
	Irish
	// stud games
	SevenStud
	SevenStudHL
	Razz
	// draw games
	FiveDraw
	Deuce7
	Deuce73Draw
	Badugi
)

type GameLimit int
const (
	FixedLimit = iota
	PotLimit
	NoLimit
)

type HandRanking int
const (
	High = iota
	LowA5
	Low27
	LowBadugi
)

// Record for specifying game blind bet amounts.
type Blinds struct {
	small uint32
	big   uint32
}

type Player interface {
	DealPrivate(Card)
	DealShared([]Card)
	DealPublic(Card)
}

// Interface for which all games must implement.
type Game interface {
	Play()
	AddPlayer(Player) error
	DealPrivate(int)
	DealShared(int)
	DealPublic(int)
}

// Record for games with community cards (Hold'em, Omaha, etc)
type CommunityGame struct {
	players    []Player
	deck       Deck
	game       GameType
	limit			 GameLimit
	dealer     int
	blinds     Blinds
	maxRaises  int
}

// Create a new community card game instance.
func NewCommunityGame(players []Player, deck Deck, game GameType, limit GameLimit, blinds Blinds, maxRaises int) *CommunityGame {
	return &CommunityGame{
		players: players,
		deck: deck,
		game: game,
		limit: limit,
		dealer: 0,
		blinds: blinds,
		maxRaises: maxRaises,
	}
}

// Hold'em/Omaha/OmahaHiLo/Omaha5/Omaha5HiLo/Courchevel/CourchevelHiLo/Irish:
//	1. Post blinds
//	2. Deal 2, 4 or 5 private cards, starting to left of big blind (2 for HE, 4 for Omaha/Irish, 5 for Omaha5/Courchevel)
//		a. Courchevel also deals 1 shared card at this point
//	3. Round of betting
//	4. Deal 3 shared cards (2 for Courchevel to complete flop)
//	5. Round of betting
//		a. Irish: discard 2 hole cards
//	6. Deal 1 shared card
//	7. Round of betting
//	8. Deal 1 shared card
//	9. Round of betting
//	10. Showdown
func (game *CommunityGame) Play() {
	game.init()
	// pre-flop
	game.dealHands()
	if game.isCourchevel() {
		game.DealShared(1)
	}
	game.betting()
	// flop
	if game.isCourchevel() {
		game.DealShared(2)
	} else {
		game.DealShared(3)
	}
	game.betting()
	if game.isIrish() {
		game.discard(2)
	}
	// turn
	game.DealShared(1)
	game.betting()
	// river
	game.DealShared(1)
	game.betting()
	game.showdown()
}

// Deal the given number of cards to each player privately (i.e. face down)
// Deals player to left of big blind first, ending with dealing to big blind.
func (game *CommunityGame) DealPrivate(cards int) {
	cnt := len(game.players)
	for i := 0; i < cards; i++ {
		for j := 1; j <= cnt; j++ {
			card := game.deck.Deal()
			p := (game.dealer + j) % cnt
			game.players[p].DealPrivate(card)
		}
	}
	game.pushHands()
}

// Deal the given number of cards to the shared community card board.
func (game *CommunityGame) DealShared(cards int) {
	dealt := make([]Card, cards)
	for i := 0; i < cards; i++ {
		dealt[i] = game.deck.Deal()
	}
	for _, player := range game.players {
		player.DealShared(dealt)
	}
	game.pushHands()
}

// Panics; community card games do not have public cards. (i.e. cards dealt to an individual player face up)
func (game *CommunityGame) DealPublic(cards int) {
	panic("cannot deal public card in community game")
}

// Prepare this game for a new round of play.
func (game *CommunityGame) init() {
	// shuffle the cards
	game.deck.Shuffle()
	// advance dealer button
	game.dealer = (game.dealer + 1) % len(game.players)
	// force the blinds to post
	game.postBlinds()
}

// Is this game a variant of Courchevel?
func (game *CommunityGame) isCourchevel() bool {
	return game.game == Courchevel || game.game == CourchevelHL
}

// Is this game a variant of Irish?
func (game *CommunityGame) isIrish() bool {
	return game.game == Irish
}

// Deal each player their private hands for a community card game.
func (game *CommunityGame) dealHands() {
	switch game.game {
	case Holdem:
		game.DealPrivate(2)
	case Omaha,OmahaHL,Irish:
		game.DealPrivate(4)
	case Omaha5,Omaha5HL,Courchevel,CourchevelHL:
		game.DealPrivate(5)
	}
}

func (game *CommunityGame) pushHands() {}
func (game *CommunityGame) betting() {}
func (game *CommunityGame) showdown() {}
func (game *CommunityGame) postBlinds() {}
func (game *CommunityGame) discard(cards int) {}

// 7Stud/7StudHiLo/Razz:
//	1. Post antes, if applicable
//	2. Deal 2 private cards and 1 public card
//	3. Player with lowest door card posts bring-in
//	4. Round of betting at lower limit, action starts to that player's left
//	5. Deal 1 public card
//	6. Round of betting at lower limit, action starts on player with strongest exposed board and continues to their left
//	7. Deal 1 public card
//	9. Round of betting at higher limit, action starts on player with strongest exposed board and continues to their left
//	10. Deal 1 public card
//	11. Round of betting at higher limit, action starts on player with strongest exposed board and continues to their left
//	12. Deal 1 private card or shared card if insufficient cards left in deck to deal to every player remaining
//	13. Round of betting at higher limit, action starts on player with strongest exposed board and continues to their left
//	14. Showdown

// 5Draw/2-7Lo/2-7TripleDrawLo/Badugi:
//	1. Post ante and blinds
//	2. Deal 4 (Badugi) or 5 (all others) private cards
//	3. For each round (1 for 5Draw/2-7Lo or 3 for Triple Draw/Badugi):
//		a. Round of betting
//		b. Draw/discard
//	4. Round of betting
//	5. Showdown

// Mixed games:
//	T - Limit 2-7 Triple Draw
//	H - Limit Hold'em
//	O - Limit Omaha Eight or Better (Hi/Lo)
//	R - Razz
//	S - Limit Seven Card Stud
//	E - Limit Stud Eight or Better (Hi/Lo)
//	H - No Limit Hold'em
//	A - Pot Limit Omaha

// Hand ranking systems:
//	1. High
//	2. A-5 Lo (California Lowball)
//	3. 2-7 Lo (Kansas City Lowball)
//	4. Badugi

