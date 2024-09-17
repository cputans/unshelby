package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentMessageContent
	err = dg.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "https://nextcloud.denike.io/s/") {
		re := regexp.MustCompile("(https://nextcloud.denike.io/s/[A-Za-z0-9]{15})")
		matches := re.FindStringSubmatch(m.Content)

		if len(matches) > 0 {
			resp, _ := http.Get(matches[0])
			if resp.StatusCode < 400 {
				body, _ := io.ReadAll(resp.Body)

				doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
				image, _ := doc.Find("a#downloadFile").Eq(0).Attr("href")

				if image != "" {
					s.ChannelMessageSend(m.ChannelID, image)
				} else {
					s.ChannelMessageSend(m.ChannelID, "NO!")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "NO!")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "NO!")
		}
	} else if strings.Contains(strings.ToLower(m.Content), "setup") {
		rand := rand.IntN(6) + 1

		if rand == 3 {
			msg := "Did somebody say setup?!?  https://denike.io/about/my-setup/desktop/"
			s.ChannelMessageSend(m.ChannelID, msg)
		}
	}

}
