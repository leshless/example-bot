package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	cl "bot/pkg/client"
)

var(
	bot *tg.Bot
	bothandler *th.BotHandler
	shutdown chan os.Signal
	waitGroup sync.WaitGroup
)

const(
	GRACEFUL_SHUTDOWN_TIME = 10 * time.Second
	OMMIT_UPDATE_TIME = 10 * time.Second
)

func closer(){
	<-shutdown
	
	log.Println("Graceful shutdown initiated. Stopping bot long polling...")
	bot.StopLongPolling()
	bothandler.Stop()

	log.Println("Waiting for all active workers to finish...")
	shutdownDone := make(chan struct{}, 1)
	go func(){
		waitGroup.Wait()
		shutdownDone <- struct{}{}
	}()

	select{
	case <- shutdownDone:
	case <- time.After(GRACEFUL_SHUTDOWN_TIME):
		log.Println("The shutdown takes too long. Some workers may be stuck. Force termination...")
	}

	log.Println("Graceful shutdown comlete.")
	os.Exit(0)
}

func runBot(){
	if len(os.Args) < 2{
		log.Fatal("A valid bot token should be included in command line arguments.")
	}

	var err error
	bot, err = tg.NewBot(os.Args[1], tg.WithDefaultDebugLogger())
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
		if update.Message != nil && time.Since(time.Unix(0, update.Message.Date)) > OMMIT_UPDATE_TIME{
			log.Println("Incoming update ommited: the valid time to handle has exceeded.")
		}

		waitGroup.Add(1)
		next(bot, update)
		waitGroup.Done()
	})

	bothandler.HandleMessage(cl.MessageHandler, th.AnyMessage())

	bothandler.Start()
}

func main() {
	shutdown = make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go closer()

	runBot()
}
