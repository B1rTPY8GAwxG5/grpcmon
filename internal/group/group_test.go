package group_test

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/group"
)

func entry(method string, code codes.Code) capture.Entry {
	return capture.Entry{
		ID:     method,
		Method: method,
		Status: status.New(code, ""),
	}
}

func TestApply_ByMethod_GroupsCorrectly(t *testing.T) {
	s := group.New(group.ByMethod)
	entries := []capture.Entry{
		entry("/svc/A", codes.OK),
		entry("/svc/B", codes.OK),
		entry("/svc/A", codes.NotFound),
	}
	groups := s.Apply(entries)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "/svc/A" {
		t.Errorf("expected first key /svc/A, got %s", groups[0].Key)
	}
	if len(groups[0].Entries) != 2 {
		t.Errorf("expected 2 entries in group A, got %d", len(groups[0].Entries))
	}
}

func TestApply_ByStatus_GroupsCorrectly(t *testing.T) {
	s := group.New(group.ByStatus)
	entries := []capture.Entry{
		entry("/svc/A", codes.OK),
		entry("/svc/B", codes.NotFound),
		entry("/svc/C", codes.OK),
	}
	groups := s.Apply(entries)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	s := group.New(nil)
	groups := s.Apply(nil)
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestNew_NilKeyFuncDefaultsToByMethod(t *testing.T) {
	s := group.New(nil)
	entries := []capture.Entry{entry("/svc/X", codes.OK)}
	groups := s.Apply(entries)
	if groups[0].Key != "/svc/X" {
		t.Errorf("unexpected key: %s", groups[0].Key)
	}
}
