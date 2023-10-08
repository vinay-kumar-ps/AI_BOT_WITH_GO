package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// "strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

func printCommandEvents(analyticschannel <-chan *slacker.CommandEvent) {
	for event := range analyticschannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition{

		Description: "send any question to wolfram",
		// Examples: "who is the president of india",
		Handler: func(botctx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")

			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data[:])

			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")

			fmt.Println(value)
			response.Reply("received")

		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
