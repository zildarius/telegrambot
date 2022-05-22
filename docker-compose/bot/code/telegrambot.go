package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zildarius/telegrambot-Go/docker-compose/bot/code/wiki"
	"github.com/zildarius/telegrambot/docker-compose/bot/code/jokes"
	"os"
	"reflect"
	"strconv"
	"time"
)

var buttonToJokes = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "/next"},
}

const wikiAnswer = 1

const jokesAnswer = 2

func telegramBot() {

	dbSwitchOn := os.Getenv("DB_SWITCH") == "on"

	//Create bot
	fmt.Println("TOKEN: " + os.Getenv("TOKEN"))
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
		//fmt.Println("panic")
	}

	bot.Debug = true

	//Set update timeout
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Get updates from bot
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		answerType := wikiAnswer // по умолчанию вики
		if dbSwitchOn {
			answerType, err = getAnswerType(update.Message.Chat.ID)
			if err != nil {
				panic(err)
			}
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			var shortAnswer = ""

			switch update.Message.Text {
			case "/start":
				shortAnswer = "Привет! Я тестовый бот. Я могу показывать шутки или искать в вики статьи."
			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {

					//Assigning number of users to num variable
					num, err := getNumberOfUsers()
					if err != nil {
						// SEND MESSAGE
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						bot.Send(msg)
						fmt.Println(err)
					}

					//Creating string which contains number of users
					shortAnswer = fmt.Sprintf("%d peoples used this bot", num)
				} else {
					shortAnswer = "Database not connected, so i can't say you how many peoples used me."
				}
			case "/wiki_search":
				fmt.Println("wiki search")
				answerType = wikiAnswer
				if dbSwitchOn {
					fmt.Println("wiki search")
					setAnswerType(update.Message.Chat.ID, wikiAnswer)
					fmt.Println("to wiki search changed")
				}
				shortAnswer = "Вы переключились на поиск по вики, для поиска введите запрос"
			case "/chuck_jokes":
				fmt.Println("Chuck jokes")
				answerType = jokesAnswer
				if dbSwitchOn {
					setAnswerType(update.Message.Chat.ID, jokesAnswer)
					fmt.Println("to chuck jokes changed")
				}
				shortAnswer = "Вы переключились в режим получения шуток. Для получения следующей шутки введите команду /next или введите номер шутки."
			case "/next":
				if answerType == wikiAnswer {
					shortAnswer = "Вы находитесь в режиме поиска по вики, для переключения на режим шуток введите /chuck_jokes"
				} else {
					fmt.Println("came to NEXT chuck jokes")
					shortAnswer = jokes.ReturnNextJoke("")
				}
			default:

				var message []string

				fmt.Println("Answer type: ", answerType)

				if answerType == wikiAnswer {

					fmt.Println("came to wiki")
					//Set search language
					language := os.Getenv("LANGUAGE")

					//Create search url
					ms, _ := wiki.URLEncoded(update.Message.Text)

					url := ms
					request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"

					//assigning value of shortAnswer slice to variable message
					message = wiki.WikipediaAPI(request)
				} else {
					fmt.Println("came to chuck jokes")
					if _, err := strconv.Atoi(update.Message.Text); err == nil {
						shortAnswer = jokes.ReturnNextJoke(update.Message.Text)
					} else {
						shortAnswer = "Введите номер шутки или выберите команду /next."
					}
				}

				if dbSwitchOn {
					//Putting username, chat_id, message, shortAnswer to database
					if err := collectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, message); err != nil {
						//Send message
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, but bot still working.")
						bot.Send(msg)
					}
				}

				//Loop throug message slice
				for _, val := range message {
					//Send message
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}

			}
			//Send message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, shortAnswer)
			if answerType == jokesAnswer {
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttonToJokes)
			}
			bot.Send(msg)
		} else {

			//Send message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите что-нибудь для поиска")
			bot.Send(msg)
		}
	}
}

func main() {

	fmt.Println("started")

	fmt.Println("CREATE_TABLE=" + os.Getenv("CREATE_TABLE"))
	fmt.Println("DB_SWITCH=" + os.Getenv("DB_SWITCH"))

	time.Sleep(10 * time.Second)

	//Creating Table
	if os.Getenv("DB_SWITCH") == "on" {

		fmt.Println("first point")

		if os.Getenv("CREATE_TABLE") == "yes" {
			fmt.Println("second point")
			if err := createTable(); err != nil {
				fmt.Println("third  point")
				panic(err)
			}
		}
	}

	time.Sleep(10 * time.Second)

	fmt.Println("starting finish!!!")

	//Call Bot
	telegramBot()
}
