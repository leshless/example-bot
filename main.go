package main

import (
	cl "bot/pkg/client"
	st "bot/pkg/storage"
	"bot/pkg/text"
	"bot/pkg/ui"
	"sync"

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
	waitgroup sync.WaitGroup
)

const(
	GRACEFUL_SHUTDOWN_TIME = 10 * time.Second
	OMIT_UPDATE_TIME = 10 * time.Second
)

func closer(){
	// listening to system call to terminate work
	<- shutdown.Done()
	
	log.Println("Graceful shutdown initiated. Stopping bot long polling...")
	bot.StopLongPolling()
	bothandler.Stop()

	log.Println("Waiting for all active workers to finish...")

	done := make(chan struct{}, 1)
	go func(){
		waitgroup.Wait()
		done <- struct{}{}
	}()

	select{
	case <-done:
		log.Println("All workers terminated successfully.")
	case <-time.After(GRACEFUL_SHUTDOWN_TIME):
		log.Println("Seems that shutdown operation takes too much time. Force termination...")
	}

	log.Println("Closing the database connection...")
	st.Close()

	log.Println("Graceful shutdown comlete.")


	os.Exit(0)
}

func load(){
	log.Println("Loading static data...")
	
	err := ui.Load()
	if err != nil{
		log.Fatal(err)
	}
	log.Println("UI data loaded sucessfuly.")

	err = text.Load()
	if err != nil{
		log.Fatal(err)
	}
	log.Println("Text data loaded sucessfuly.")

	log.Println("Initializing DB connection...")
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

	cl.SetBot(bot)

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatal(err)
	}

	bothandler, err = th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}
	
	bothandler.Use(func(bot *tg.Bot, update tg.Update, next th.Handler) {
		// Check if update is not too old to be processed.
		if update.Message != nil{
			updatetime := time.Unix(update.Message.Date, 0)
			if time.Since(updatetime) > OMIT_UPDATE_TIME{
				log.Println("Incoming update omitted: a valid response time has elapsed.")
				return
			}
		}

		waitgroup.Add(1)
		next(bot, update)
		waitgroup.Done()
	})

	bothandler.HandleCallbackQuery(cl.QueryHandler, th.AnyCallbackQuery())
	bothandler.HandleMessage(cl.CommandHandler, th.AnyCommand())
	bothandler.HandleMessage(cl.MessageHandler, th.AnyMessage())

	log.Print("Bot is running.")
	bothandler.Start()
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
