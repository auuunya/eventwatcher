package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/auuunya/eventwatcher"
)

func main() {
	ctx := context.TODO()
	notify := eventwatcher.NewEventNotifier(ctx)
	defer notify.Close()

	channels := []string{"Application", "System", "Microsoft-Windows-Kernel-Dump/Operational"}
	for _, channel := range channels {
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
