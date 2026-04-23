package transform_test

import (
	"testing"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/transform"
)

func TestRedactMetadataKey_RedactsMatchingKey(t *testing.T) {
	e := baseEntry()
	e.Metadata = map[string]string{"Authorization": "Bearer secret", "x-id": "123"}

	f := transform.RedactMetadataKey("authorization")
	out, ok := f(e)
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if out.Metadata["Authorization"] != "REDACTED" {
		t.Errorf("expected REDACTED, got %q", out.Metadata["Authorization"])
	}
	if out.Metadata["x-id"] != "123" {
		t.Errorf("unrelated key should be unchanged, got %q", out.Metadata["x-id"])
	}
}

func TestRedactMetadataKey_NoMetadata_ReturnsUnchanged(t *testing.T) {
	e := baseEntry()
	f := transform.RedactMetadataKey("authorization")
	out, ok := f(e)
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if out.Metadata != nil {
		t.Error("expected nil metadata")
	}
}

func TestKeepMethods_AllowsListedMethod(t *testing.T) {
	f := transform.KeepMethods("/svc/Method", "/svc/Other")
	_, ok := f(baseEntry())
	if !ok {
		t.Fatal("expected entry to be kept")
	}
}

func TestKeepMethods_DropsUnlistedMethod(t *testing.T) {
	f := transform.KeepMethods("/svc/Other")
	_, ok := f(baseEntry())
	if ok {
		t.Fatal("expected entry to be dropped")
	}
}

func TestNormaliseMethod_TrimsAndLowercases(t *testing.T) {
	e := capture.Entry{Method: "//Svc/MyMethod"}
	f := transform.NormaliseMethod()
	out, ok := f(e)
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if out.Method != "svc/mymethod" {
		t.Errorf("expected svc/mymethod, got %q", out.Method)
	}
}

func TestPreset_ChainedRedactAndKeep(t *testing.T) {
	entries := []capture.Entry{
		{ID: "1", Method: "/svc/A", Metadata: map[string]string{"token": "abc"}},
		{ID: "2", Method: "/svc/B", Metadata: map[string]string{"token": "xyz"}},
	}

	c := transform.New().
		Add(transform.KeepMethods("/svc/A")).
		Add(transform.RedactMetadataKey("token"))

	out := c.ApplyAll(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Metadata["token"] != "REDACTED" {
		t.Errorf("expected REDACTED, got %q", out[0].Metadata["token"])
	}
}
