package util

import "math/rand"

// List of kramerisms
var kramerisms = []string{
	"Well, you’re just as pretty as any of them. You just need a nose job.",
	"Who's gonna turn down a Junior Mint? It's chocolate, it's peppermint, it's delicious.",
	"I’m out there, Jerry. And I’m loving every minute of it!",
	"I love the name ‘Isosceles.’ If I had a kid, I would name him Isosceles. Isosceles Kramer.",
	"Whadda you think junior? These hands have been soaking in ivory liquid?",
	"Now what does the little man inside you say? See you gotta listen to the little man.",
	"You contribute nothing to society!",
	"Marriage? Family? They're prisons...Man-made prisons. You're doing time...",
	"You know what would make a great coffee table book? A coffee table book about coffee tables! Get it?",
	"Because I'm like ice, buddy. When I don't like you, you've got problems.",
	"Oh, Jerry, wake up to reality. It's military thing. They're probably creating a whole army of pig warriors.",
	"Hey, Silvio, look at Jerry here, prancing around in his coat with his purse. Yup, he's a dandy. He's a real fancy boy.",
}

// Function to randomly select a kramerism
func GetKramerism() string {
	return kramerisms[rand.Intn(len(kramerisms))]
}
