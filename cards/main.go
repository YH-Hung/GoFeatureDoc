package main

var card string = "latter assign"

func main() {
	//cards := newDeck()
	//
	//hand, remainingCards := deal(cards, 5)
	//hand.print()
	//fmt.Println("-----")
	//remainingCards.print()
	//
	//cards.saveToFile("my_cards")

	cards := newDeckFromFile("my_cards")
	cards.shuffle()
	cards.print()
}
