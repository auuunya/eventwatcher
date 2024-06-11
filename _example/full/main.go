package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/auuunya/eventwatcher"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	session, err := eventwatcher.EvtOpenChannelEnum(0)
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		return
	}
	defer eventwatcher.EvtClose(session)
	channels, err := eventwatcher.EvtNextChannelPath(session)
	if err != nil {
		return
	}

	ctx := context.TODO()
	notify := eventwatcher.NewEventNotifier(ctx)
	defer notify.Close()

	for _, channel := range channels[:100] {
		fmt.Printf("channel: %v\n", channel)
		err := notify.AddWatcher(channel)
		if err != nil {
			continue
		}
	}
	go func() {
		for ch := range notify.EventLogChannel {
			e := eventwatcher.ParseEventLogData(ch)
			fmt.Printf("e: %+v\n", e)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutting down\n")
}
