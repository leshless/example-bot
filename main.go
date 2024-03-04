package main

import (
	cl "bot/pkg/client"
	st "bot/pkg/storage"
	"bot/pkg/ui"

	ctx "context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

var(
	bot *tg.Bot
	bothandler *th.BotHandler
	shutdown ctx.Context
)

const(
	GRACEFUL_SHUTDOWN_TIME = 10 * time.Second
	OMMIT_UPDATE_TIME = 10 * time.Second
)

func closer(){
	// listening to system call to terminate work
	<- shutdown.Done()
	
	log.Println("Graceful shutdown initiated. Stopping bot long polling...")
	bot.StopLongPolling()
	bothandler.Stop()

	log.Println("Waiting for all active workers to finish...")

	//

	log.Println("Graceful shutdown comlete.")

	st.Close()

	os.Exit(0)
}

func load(){
	log.Println("Loading static data...")
	
	err := ui.Load()
	if err != nil{
		log.Fatal(err)
	}
	log.Println("UI data loaded sucessfuly.")

	log.Println("Initializing database connection...")
	err = st.Init()
	if err != nil{
		log.Fatal(err)
	}
	log.Println("DB connected sucessfuly.")

	log.Println("Loading client resources...")
	err = cl.Load()
	if err != nil{
		log.Fatal(err)
	}

	log.Println("Resources loaded sucessfuly.")
	//
}

func runBot(){
	log.Print("Starting the bot...")

	var err error
	bot, err = tg.NewBot(os.Args[1], tg.WithDefaultLogger(false, true))
	if err != nil {
		log.Fatal(err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatal(err)
	}

	bothandler, err = th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}
	
	bothandler.Use(func(bot *tg.Bot, update tg.Update, next th.Handler) {
		// // wrapper
		// done := make(chan struct{}, 1)
		// go func(){
		// 	next(bot, update)
		// 	done <- struct{}{}
		// }()

		// select{
		// case <- done:
		// 	// everything is fine
		// case <- context.Done():
		// 	// imediately stop work and perform some auxilary purge
		// }
		
		next(bot, update)
	})


	bothandler.Start()
	log.Print("Bot is ")
}

func main() {
	if len(os.Args) < 2{
		log.Fatal("A valid bot token should be included in command line arguments.")
	}

	shutdown, _ = signal.NotifyContext(ctx.Background(), syscall.SIGINT, syscall.SIGTERM)
	go closer()
	
	load()
	runBot()
} 
