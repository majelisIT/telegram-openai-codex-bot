package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {

	log.Println("starting bot")

	api := NewCodexApi(os.Getenv("OPENAI_AI_KEY"))

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/ask", func(c tele.Context) error {

		var answer string

		question := c.Message().Payload

		result, err := api.GetCodexSuggestion(question)

		if len(result.Choices) > 0 {
			answer = result.Choices[0].Text
		} else {
			answer = "I don't know"
		}

		if err != nil {
			answer = err.Error()
		}

		answer += "\n\n generated using text-davinci-003 model"

		return c.Reply(answer)
	})

	b.Start()

}
