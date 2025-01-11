package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Card struct {
	Suit  string
	Value int
}

type Player struct {
	Name  string
	Hand  []Card
	Score int
	Pile  []Card
}

var suits = []string{"Hearts", "Diamonds", "Clubs", "Spades"}
var values = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} // 1=Ace, 11=Jack, 12=Queen, 13=King

func createDeck() []Card {
	var deck []Card
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

// Uses local time to generate random seed
// for truly random shuffling
func shuffleDeck(deck []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func dealCards(deck []Card, numCards int) ([]Card, []Card) {
	return deck[:numCards], deck[numCards:]
}

func playCard(player *Player, index int) Card {
	card := player.Hand[index]
	player.Hand = append(player.Hand[:index], player.Hand[index+1:]...)
	return card
}

func calculateScore(player *Player) {
	player.Score = 0
	for _, card := range player.Pile {
		if card.Value == 1 {
			player.Score += 1 // Ace is worth 1 point
		} else if card.Value == 11 {
			player.Score += 1 // Jack is worth 1 point
		} else if card.Value == 10 && card.Suit == "Diamonds" {
			player.Score += 3 // 10 of Diamonds is worth 3 points
		} else if card.Value == 2 && card.Suit == "Clubs" {
			player.Score += 2 // 2 of Clubs is worth 2 points
		}
	}
}

func createCardString(card Card) string {
	return fmt.Sprintf("%d of %s", card.Value, card.Suit)
}

func main() {
	app := tview.NewApplication()

	deck := createDeck()
	deck = shuffleDeck(deck)

	// Initialize players
	player1 := Player{Name: "Player 1"}
	player2 := Player{Name: "Player 2"}
	players := []*Player{&player1, &player2}

	// Initial deal
	player1.Hand, deck = dealCards(deck, 4)
	player2.Hand, deck = dealCards(deck, 4)
	table, deck := dealCards(deck, 4)

	playerTurn := 0

	// Layout setup
	handText := tview.NewTextView().SetText(fmt.Sprintf("%s's Hand: %v", players[playerTurn].Name, player1.Hand))
	tableText := tview.NewTextView().SetText(fmt.Sprintf("Table: %v", table))
	statusText := tview.NewTextView().SetText(fmt.Sprintf("Player turn: %s", players[playerTurn].Name))

	// Input field for card index
	cardIndexInput := tview.NewInputField().SetLabel("Card Index: ").SetFieldWidth(10)

	// Update function to refresh the UI
	updateUI := func() {
		handText.SetText(fmt.Sprintf("%s's Hand: %v", players[playerTurn].Name, players[playerTurn].Hand))
		tableText.SetText(fmt.Sprintf("Table: %v", table))
		statusText.SetText(fmt.Sprintf("Player turn: %s", players[playerTurn].Name))
	}

	// Handle card play logic
	cardIndexInput.SetDoneFunc(func(key tcell.Key) {
		index := 0
		_, _ = fmt.Sscanf(cardIndexInput.GetText(), "%d", &index)
		if index < 1 || index > len(players[playerTurn].Hand) {
			statusText.SetText("Invalid card selection")
			return
		}

		playedCard := playCard(players[playerTurn], index-1)
		table = append(table, playedCard)

		// Check for capturing
		var newTable []Card
		for _, tableCard := range table {
			if tableCard.Value == playedCard.Value {
				players[playerTurn].Pile = append(players[playerTurn].Pile, tableCard)
			} else {
				newTable = append(newTable, tableCard)
			}
		}
		table = newTable

		// Switch turn
		playerTurn = (playerTurn + 1) % 2
		updateUI()
	})

	// Layout for the interface
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Basra Game").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(handText, 0, 1, false).
		AddItem(tableText, 0, 1, false).
		AddItem(cardIndexInput, 1, 0, true).
		AddItem(statusText, 1, 0, false)

	// Start the app
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
