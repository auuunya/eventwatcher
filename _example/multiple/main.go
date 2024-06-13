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
		err := notify.AddWatcher(channel)
		if err != nil {
			continue
		}
	}

	go func() {
		for ch := range notify.EventLogChannel {
			record := eventwatcher.ParseEventLogData(ch.Buffer)
			fmt.Printf("name: %s, handle: %v, record: %+v\n", ch.Name, ch.Handle, record)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutting down\n")
}
