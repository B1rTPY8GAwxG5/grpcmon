package compare_test

import (
	"bytes"
	"testing"
	"time"

	"grpcmon/internal/capture"
	"grpcmon/internal/compare"
	google_grpc_codes "google.golang.org/grpc/codes"
)

func makeStore(entries []capture.Entry) *capture.Store {
	s := capture.NewStore(10)
	for _, e := range entries {
		s.Add(e)
	}
	return s
}

func entry(method string, code google_grpc_codes.Code, resp string) capture.Entry {
	return capture.Entry{
		ID:       "id-" + method,
		Method:   method,
		Status:   code,
		Response: resp,
		Latency:  10 * time.Millisecond,
	}
}

func TestStores_MatchingEntries(t *testing.T) {
	base := makeStore([]capture.Entry{entry("/svc/Hello", 0, `{"ok":true}`)})
	cand := makeStore([]capture.Entry{entry("/svc/Hello", 0, `{"ok":true}`)})

	r := compare.Stores(base, cand)
	if r.MatchCount != 1 || r.MismatchCount != 0 {
		t.Fatalf("expected 1 match, got %+v", r)
	}
}

func TestStores_MismatchedResponse(t *testing.T) {
	base := makeStore([]capture.Entry{entry("/svc/Hello", 0, `{"ok":true}`)})
	cand := makeStore([]capture.Entry{entry("/svc/Hello", 0, `{"ok":false}`)})

	r := compare.Stores(base, cand)
	if r.MismatchCount != 1 {
		t.Fatalf("expected 1 mismatch, got %+v", r)
	}
}

func TestStores_SkipsMissingCandidate(t *testing.T) {
	base := makeStore([]capture.Entry{entry("/svc/Missing", 0, "")})
	cand := makeStore([]capture.Entry{})

	r := compare.Stores(base, cand)
	if len(r.Results) != 0 {
		t.Fatalf("expected no results, got %d", len(r.Results))
	}
}

func TestFprint_ContainsSummary(t *testing.T) {
	base := makeStore([]capture.Entry{entry("/svc/A", 0, "x")})
	cand := makeStore([]capture.Entry{entry("/svc/A", 0, "x")})
	r := compare.Stores(base, cand)

	var buf bytes.Buffer
	compare.Fprint(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("Compared")) {
		t.Fatalf("expected summary line, got: %s", buf.String())
	}
}
