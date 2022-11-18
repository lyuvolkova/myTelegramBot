package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

func getUserId(str string) (map[int64]struct{}, error) {
	parts := strings.Split(str, ",")
	users := make(map[int64]struct{})
	for _, val := range parts {
		userId, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		users[userId] = struct{}{}
	}
	return users, nil
}

func main() {
	coldData := ""
	users, err := getUserId(os.Getenv("USER_IDS"))
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	myCron := cron.New(cron.WithParser(secondParser), cron.WithChain())
	myCron.Start()
	defer myCron.Stop()
	myCron.AddFunc("0 00 19 18 * ?", func() {
		for i, _ := range users {
			_, err = bot.Send(tgbotapi.NewMessage(i, "Hello\nhow are u?\nI NEED DATA!"))
			if err != nil {
				log.Println(fmt.Errorf("cron: %w", err))
			}
		}

	})
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if _, exists := users[update.Message.Chat.ID]; !exists { //check the sender of a message
				continue
			}
			text := update.Message.Text
			log.Printf("[%s] %s\n", update.Message.From.UserName, text)
			textMsg := "I don't understand you !"
			if text == "Hello" {
				textMsg = "Hello, " + update.Message.Chat.FirstName + "!"
			} else if text == "Bye" {
				textMsg = "Bye bye!"
			} else if text == "/cold_data" {
				textMsg = "Cold water: " + getPrevData(coldData)
			} else if strings.HasPrefix(text, "/cold_data ") {
				textMsg = writeColdData(text)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
			//msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func writeColdData(text string) string {
	textMsg := "Ok, date saved"
	coldData := text[11:] //string([]rune(text)[11:])
	if _, err := strconv.ParseFloat(coldData, 32); err != nil {
		textMsg = "Incorrect number"
	} else {
		if !checkData(coldData) {
			textMsg = "Check number: new data < prev data"
		} else {
			dt := time.Now()
			coldData += dt.Format("  (01.02.2006 15:04:05)") + "\n"
			file, err := os.OpenFile("coldData.txt", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			if _, err = file.WriteString(coldData); err != nil {
				log.Fatal(err)
			}
		}
	}
	return textMsg
}

func checkData(stringDate string) bool {
	var d2 string
	d2 = getPrevData(d2)
	d1, err := strconv.ParseFloat(stringDate, 32)
	if err != nil {
		log.Fatal(err)
	}
	countNum := 0
	for i, _ := range d2 {
		if string(d2[i]) == " " {
			countNum = i
			break
		}
	}
	newD2, err := strconv.ParseFloat(d2[:countNum], 32)
	if d1 < newD2 {
		return false
	}
	return true
}
func getPrevData(coldData string) string {
	file, err := excelize.OpenFile("WaterData.xlsx")
	if err != nil {
		log.Println(err)
		return "Not open file"
	}
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := file.GetCols("Water data")
	if err != nil {
		log.Fatal(err)
	}
	coldData = rows[1][len(rows[1])-1]
	if coldData == "" || coldData == "ХВС" {
		coldData = "No previous data"
	}
	return coldData
}
