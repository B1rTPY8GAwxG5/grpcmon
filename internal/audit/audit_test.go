package audit

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestRecord_And_List(t *testing.T) {
	l := New(10)
	l.Record(KindReplay, "/svc/Method", nil)
	l.Record(KindExport, "out.json", nil)

	events := l.List()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Kind != KindReplay {
		t.Errorf("expected first kind %q, got %q", KindReplay, events[0].Kind)
	}
	if events[1].Detail != "out.json" {
		t.Errorf("unexpected detail: %s", events[1].Detail)
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	l := New(3)
	for i := 0; i < 4; i++ {
		l.Record(KindFilter, string(rune('a'+i)), nil)
	}
	events := l.List()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0].Detail != "b" {
		t.Errorf("expected oldest to be evicted, got detail %q", events[0].Detail)
	}
}

func TestRecord_StoresError(t *testing.T) {
	l := New(5)
	sentinel := errors.New("connection refused")
	l.Record(KindReplay, "detail", sentinel)

	events := l.List()
	if events[0].Err != sentinel {
		t.Errorf("expected stored error, got %v", events[0].Err)
	}
}

func TestClear_RemovesAll(t *testing.T) {
	l := New(10)
	l.Record(KindSnapshot, "snap1", nil)
	l.Clear()
	if len(l.List()) != 0 {
		t.Error("expected empty log after Clear")
	}
}

func TestNew_DefaultMax(t *testing.T) {
	l := New(0)
	if l.max != 256 {
		t.Errorf("expected default max 256, got %d", l.max)
	}
}

func TestRecord_ConcurrentSafety(t *testing.T) {
	l := New(512)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Record(KindReplay, "concurrent", nil)
		}()
	}
	wg.Wait()
	if len(l.List()) != 50 {
		t.Errorf("expected 50 events, got %d", len(l.List()))
	}
}

func TestFprint_ContainsEvents(t *testing.T) {
	l := New(10)
	l.Record(KindExport, "result.json", nil)

	var buf bytes.Buffer
	Fprint(&buf, l)

	out := buf.String()
	if !strings.Contains(out, "export") {
		t.Errorf("expected 'export' in output, got: %s", out)
	}
	if !strings.Contains(out, "result.json") {
		t.Errorf("expected detail in output, got: %s", out)
	}
}

func TestEvent_String_ErrorFormat(t *testing.T) {
	e := Event{Kind: KindReplay, Detail: "/pkg/Method", Err: errors.New("timeout")}
	s := e.String()
	if !strings.Contains(s, "err: timeout") {
		t.Errorf("expected error in string, got: %s", s)
	}
}
