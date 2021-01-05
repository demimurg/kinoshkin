package main

import (
	"kinoshkin/internal/conferencier"
	bot "kinoshkin/internal/telebot"
)

func main() {
	bot.New(conferencier.Mock{}).Start()
}
