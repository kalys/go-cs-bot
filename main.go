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

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// user := tb.User{ID: 62540640}
	// b.Send(&user, "aoeuaoeuaoeua aoeuaoeu")
	// return

	b.Handle("/players", func(m *tb.Message) {
		address := fmt.Sprintf("%s:27015", cs_host)
		// address = "46.42.18.18:27018"
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

		fmt.Println(m.Sender)

		response := fmt.Sprintf("```\nPlayers: %d\nBots: %d\nMap: %s%s```", info.Players, info.Bots, info.MapName, buffer.String())

		b.Send(m.Chat, response, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	b.Start()
}
