package eventwatcher

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestMemSpike(t *testing.T) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	t.Logf("Before: Alloc=%d TotalAlloc=%d Sys=%d NumGC=%d", m1.Alloc, m1.TotalAlloc, m1.Sys, m1.NumGC)

	// Create notifier and watcher
	ctx := context.Background()
	n := NewEventNotifier(ctx)
	defer n.Close()

	f, err := os.CreateTemp("", "mem_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	if err := n.AddWatcher(f.Name()); err != nil {
		t.Fatalf("failed to add watcher: %v", err)
	}

	// wait a bit for goroutines to start
	time.Sleep(200 * time.Millisecond)

	runtime.GC()
	runtime.ReadMemStats(&m2)
	t.Logf("After: Alloc=%d TotalAlloc=%d Sys=%d NumGC=%d", m2.Alloc, m2.TotalAlloc, m2.Sys, m2.NumGC)

	if m2.Alloc > m1.Alloc*10+1024*1024 {
		t.Logf("Large allocation detected: before=%d after=%d", m1.Alloc, m2.Alloc)
	}
}