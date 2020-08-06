package main

import (
	bot "bot/bot"
	"os"
)

func main() {
	os.Remove(os.Getenv("APP_OUTPUT_FILE"))
	bot.Bot()
}
