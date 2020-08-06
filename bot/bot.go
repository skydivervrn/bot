package bot

import (
	"bot/settings"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	appOut   = os.Getenv("APP_OUTPUT_FILE")
	cmdOut   = os.Getenv("CMD_OUTPUT_FILE")
	answer   []string
	answerP  = &answer
	counter  = 1
	counterP = &counter
	update   tgbotapi.Update
	updP     = &update
	botP     *tgbotapi.BotAPI
)

// Bot some comment
func Bot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	botP = bot
	botP.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := botP.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		*updP = update
		receiver()
	}
}

func receiver() {
	if !accessCheck() {
		botP.Send(message(fmt.Sprintf("You do not have an access, username: %s", updP.Message.Chat.UserName)))
		return
	}
	ctmp := *counterP
	if !writeFile(strings.Join([]string{strconv.Itoa(*counterP), updP.Message.Text}, "-\n")) {
		log.Printf("DEBUG: %v", strconv.Itoa(*counterP))
		*counterP = *counterP + 1
	}
	errR := readFile()
	if errR != nil {
		log.Printf("File reading error: %s", errR)
		botP.Send(message(fmt.Sprintf("File reading error: %s", errR)))
		return
	}
	time.Sleep(2 * time.Second)
	tmp := time.Now()
	for strconv.Itoa(ctmp) != (*answerP)[0] {
		time.Sleep(5 * time.Second)
		readFile()
		botP.Send(message(fmt.Sprintf("Waiting another command: %f sec", time.Now().Sub(tmp).Seconds())))
	}
	botP.Send(message(fmt.Sprintf("Result is: %s \nError is %s\n", (*answerP)[1], (*answerP)[2])))
}

func accessCheck() bool {
	if stringInSlice(updP.Message.From.UserName, settings.AdminLists.ProductionAdmins) {
		return true
	}
	return false
}

func message(text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(updP.Message.Chat.ID, text)
	msg.ReplyToMessageID = updP.Message.MessageID
	return msg
}

func writeFile(body string) bool {
	file, err := os.Create(appOut)
	if err != nil {
		log.Println("Unable to create file:", err)
		return true
	}
	_, err2 := file.WriteString(body)
	if err2 != nil {
		log.Println("Unable to write file:", err)
		return true
	}
	defer file.Close()
	return false
}

func readFile() error {
	log.Printf("Start reading answer from file...")
	dat, err := ioutil.ReadFile(cmdOut)
	if err != nil {
		log.Println("Unable to open file:", err)
		return err
	}
	*answerP = strings.Split(string(dat), "-\n")
	log.Printf("File read")
	return err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
