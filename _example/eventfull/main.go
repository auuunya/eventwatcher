package main

import (
	"bytes"
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
		return
	}
	defer eventwatcher.EvtClose(session)
	channels, err := eventwatcher.EvtNextChannelPath(session)
	if err != nil {
		return
	}

	f, err := os.OpenFile("event.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	ctx := context.Background()
	notify := eventwatcher.NewEventNotifier(ctx)
	defer notify.Close()
	for _, channel := range channels[:] {
		err := notify.AddWatcher(channel)
		if err != nil {
			continue
		}
	}

	go receiver(f, notify.EventLogChannel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutting down\n")
}

func receiver(f *os.File, channel chan *eventwatcher.EventEntry) {
	for ch := range channel {
		var buf bytes.Buffer
		buffer := ch.Buffer
		r := eventwatcher.ParseEventLogData(buffer)
		content := eventwatcher.FormatContent(buffer)
		data := fmt.Sprintf("Name: %s, Handle: %v, Length: %d,Reserved: %d,RecordNumber: %d,TimeGenerated: %d,TimeWritten: %d,EventID: %d,EventType: %d,NumStrings: %d,EventCategory: %d,ReservedFlags: %d,ClosingRecordNumber: %d,StringOffset: %d,UserSidLength: %d,UserSidOffset: %d,DataLength: %d,DataOffset: %d,Content: %s.\n",
			ch.Name,
			ch.Handle,
			r.Length,
			r.Reserved,
			r.RecordNumber,
			r.TimeGenerated,
			r.TimeWritten,
			r.EventID,
			r.EventType,
			r.NumStrings,
			r.EventCategory,
			r.ReservedFlags,
			r.ClosingRecordNumber,
			r.StringOffset,
			r.UserSidLength,
			r.UserSidOffset,
			r.DataLength,
			r.DataOffset,
			content,
		)
		buf.Write([]byte(data))
		f.Write(buf.Bytes())
	}
}
