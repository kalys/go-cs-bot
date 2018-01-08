package main

import (
	"fmt"
	// "github.com/alliedmodders/blaster/valve"
	"../blaster/valve"
	"bytes"
	"github.com/olekukonko/tablewriter"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"strconv"
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
	address := fmt.Sprintf("%s:27015", cs_host)
	// address = "46.42.18.18:27018"

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		var prevPlayers []valve.Player
		c := time.Tick(10 * time.Second)
		for now := range c {
			_ = now
			client, err := valve.NewServerQuerier(address, 10*time.Second)
			defer client.Close()
			if err != nil {
				continue
			}

			playersInfo, err := client.QueryPlayers()
			if err != nil {
				continue
			}

			new, left := playersInfo.Diff(prevPlayers)

			prevPlayers = playersInfo.Players

			var buffer bytes.Buffer
			if len(new) > 0 {
				for index, player := range new {
					buffer.WriteString(player.Nickname)
					if index < len(new)-1 {
						buffer.WriteString(", ")
					}
				}
				buffer.WriteString(" connected")
			}

			if len(left) > 0 {
				if len(new) > 0 {
					buffer.WriteString("\n")
				}

				for index, player := range left {
					buffer.WriteString(player.Nickname)
					if index < len(left)-1 {
						buffer.WriteString(", ")
					}
				}
				buffer.WriteString(" left")

			}

			if len(new) > 0 || len(left) > 0 {
				chat := tb.Chat{ID: -38942708}
				b.Send(&chat, buffer.String())
			}

		}
	}()

	b.Handle("/players", func(m *tb.Message) {
		client, err := valve.NewServerQuerier(address, 10*time.Second)
		defer client.Close()

		if err != nil {
			b.Send(m.Sender, "something wrong")
			return
		}
		info, err := client.QueryInfo()
		if err != nil {
			b.Send(m.Sender, "something wrong query info")
		}

		players, err := client.QueryPlayers()

		var buffer bytes.Buffer

		if err == nil && len(players.Players) > 0 {
			buffer.WriteString("\n\n")
			table := tablewriter.NewWriter(&buffer)
			table.SetHeader([]string{"Nickname", "Kills"})
			for _, player := range players.Players {
				duration := time.Duration(player.Time) * time.Second
				rowData := []string{player.Nickname, strconv.Itoa(int(player.Kills)), duration.String()}
				table.Append(rowData)
			}
			table.Render()
		}

		response := fmt.Sprintf("```\nPlayers: %d\nBots: %d\nMap: %s%s```", info.Players, info.Bots, info.MapName, buffer.String())

		b.Send(m.Chat, response, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	b.Start()
}
