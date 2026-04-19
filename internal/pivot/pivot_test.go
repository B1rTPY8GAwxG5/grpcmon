package pivot_test

import (
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"grpcmon/internal/capture"
	"grpcmon/internal/pivot"
)

func entry(method string, status codes.Code, ms int) capture.Entry {
	return capture.Entry{
		Method:   method,
		Status:   status,
		Duration: time.Duration(ms) * time.Millisecond,
	}
}

func TestBuild_ByMethod_GroupsCorrectly(t *testing.T) {
	entries := []capture.Entry{
		entry("/svc/Foo", codes.OK, 10),
		entry("/svc/Foo", codes.OK, 20),
		entry("/svc/Bar", codes.OK, 5),
	}
	table := pivot.Build(entries, pivot.ByMethod)
	if len(table) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(table))
	}
	if table[0].Key != "/svc/Bar" {
		t.Errorf("expected first row /svc/Bar, got %s", table[0].Key)
	}
	if table[1].Count != 2 {
		t.Errorf("expected count 2 for /svc/Foo, got %d", table[1].Count)
	}
	if table[1].AvgMS != 15 {
		t.Errorf("expected avgMS 15, got %f", table[1].AvgMS)
	}
}

func TestBuild_ByStatus_GroupsCorrectly(t *testing.T) {
	entries := []capture.Entry{
		entry("/svc/A", codes.OK, 10),
		entry("/svc/B", codes.NotFound, 20),
		entry("/svc/C", codes.NotFound, 30),
	}
	table := pivot.Build(entries, pivot.ByStatus)
	if len(table) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(table))
	}
	for _, r := range table {
		if r.Key == "NotFound" && r.Count != 2 {
			t.Errorf("expected 2 NotFound entries, got %d", r.Count)
		}
	}
}

func TestBuild_ErrorCount(t *testing.T) {
	entries := []capture.Entry{
		entry("/svc/X", codes.OK, 5),
		entry("/svc/X", codes.Internal, 15),
	}
	table := pivot.Build(entries, pivot.ByMethod)
	if len(table) != 1 {
		t.Fatalf("expected 1 row, got %d", len(table))
	}
	if table[0].ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", table[0].ErrorCount)
	}
}

func TestBuild_EmptyEntries(t *testing.T) {
	table := pivot.Build(nil, pivot.ByMethod)
	if len(table) != 0 {
		t.Errorf("expected empty table, got %d rows", len(table))
	}
}
