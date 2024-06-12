package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/auuunya/eventwatcher"
)

func main() {
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

	for _, channel := range channels[:] {
		err := notify.AddWatcher(channel)
		if err != nil {
			continue
		}
	}
	go func() {
		for ch := range notify.EventLogChannel {
			e := eventwatcher.ParseEventLogData(ch)
			fmt.Printf("event log changed: %#v\n", e.RecordNumber)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutting down\n")
}
