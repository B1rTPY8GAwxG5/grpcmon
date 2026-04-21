package mask_test

import (
	"testing"

	"github.com/user/grpcmon/internal/capture"
	"github.com/user/grpcmon/internal/mask"
)

func entry(meta map[string]string) capture.Entry {
	return capture.Entry{
		Method:   "/svc/Method",
		Metadata: meta,
	}
}

func TestApply_RedactsMatchingField(t *testing.T) {
	m := mask.New("authorization")
	e := entry(map[string]string{"authorization": "Bearer secret", "x-request-id": "abc"})
	got := m.Apply(e)
	if got.Metadata["authorization"] != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %q", got.Metadata["authorization"])
	}
	if got.Metadata["x-request-id"] != "abc" {
		t.Fatalf("unmasked field should be unchanged")
	}
}

func TestApply_CaseInsensitiveKey(t *testing.T) {
	m := mask.New("Authorization")
	e := entry(map[string]string{"AUTHORIZATION": "token"})
	got := m.Apply(e)
	if got.Metadata["AUTHORIZATION"] != "[REDACTED]" {
		t.Fatalf("expected case-insensitive match to redact value")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m := mask.New("secret")
	orig := entry(map[string]string{"secret": "mysecret"})
	_ = m.Apply(orig)
	if orig.Metadata["secret"] != "mysecret" {
		t.Fatal("Apply must not mutate the original entry")
	}
}

func TestApply_NoMetadata_ReturnsUnchanged(t *testing.T) {
	m := mask.New("authorization")
	e := capture.Entry{Method: "/svc/Ping"}
	got := m.Apply(e)
	if got.Method != e.Method {
		t.Fatal("entry with no metadata should be returned unchanged")
	}
}

func TestApplyAll_MasksEveryEntry(t *testing.T) {
	m := mask.New("token")
	entries := []capture.Entry{
		entry(map[string]string{"token": "abc"}),
		entry(map[string]string{"token": "xyz", "other": "keep"}),
	}
	got := m.ApplyAll(entries)
	for i, e := range got {
		if e.Metadata["token"] != "[REDACTED]" {
			t.Fatalf("entry %d: expected [REDACTED]", i)
		}
	}
	if got[1].Metadata["other"] != "keep" {
		t.Fatal("non-masked field should be preserved")
	}
}
