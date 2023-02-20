package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {

	log.Println("starting bot")

	api := NewCodexApi(os.Getenv("OPENAI_AI_KEY"))

	codexApi := NewDalleApi(os.Getenv("OPENAI_AI_KEY"))

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/info", func(c tele.Context) error {
		return c.Reply("Majelis IT telenot openai v1.1.0 Full commits https://github.com/majelisIT/telegram-openai-codex-bot/releases/tag/v1.1.0")
	})

	b.Handle("/ask", func(c tele.Context) error {

		if !c.Message().FromGroup() {
			return c.Reply("Not Allowed, invite to group first")
		}

		msg, err := b.Reply(c.Message(), "Okay tunggu 🙏 ...")
		if err != nil {
			return c.Reply("Gagal coba lagi")
		}

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
		_, err = b.Edit(msg, answer)

		return err
	})

	b.Handle("/image", func(c tele.Context) error {

		msg, err := b.Reply(c.Message(), "Okay wait ...")
		if err != nil {
			return c.Reply("Gagal coba lagi")
		}
		time.Sleep(time.Second)

		command := c.Message().Payload

		b.Edit(msg, fmt.Sprintf("Mencoba membuat gambar %s 🫡🫡🫡🫡🫡", command))

		res, err := codexApi.GetDalleImage(command)
		if err != nil {
			return c.Reply("Maaf sedang gak bisa bikin gambar, coba lagi deh")
		}

		b.Edit(msg, "Gambar berhasil dibuat 🍻")

		imageResult := &tele.Photo{File: tele.FromURL(res.Data[0].URL)}

		return c.Reply(imageResult)
	})

	b.Start()

}
