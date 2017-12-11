package main

import (
	"fmt"
	"github.com/kalys/SteamCondenserGo/servers"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)

func main() {
	token, ok := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !ok {
		fmt.Printf("%s not set\n", "TELEGRAM_BOT_TOKEN")
	}
	cs_host, ok := os.LookupEnv("CS_HOST")
	if !ok {
		fmt.Printf("%s not set\n", "CS_HOST")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	b.Handle("/players", func(m *tb.Message) {
		address := fmt.Sprintf("%s:27015", cs_host)
		goldServer := servers.GoldServer{
			Address: address,
		}
		info, err := goldServer.GetInfo()

		if err != nil {
			b.Send(m.Sender, "something wrong")
		}

		response := fmt.Sprintf("Players: %d\nBots: %d", info.NumPlayers, info.Bots)

		b.Send(m.Chat, response)
	})

	b.Start()
}
