package main

import "testing"

func TestEncodeKnownValues(t *testing.T) {
	tests := []struct {
		id   int64
		want string
	}{
		{0, "AAAA"},
		{1, "AAAQ"},
		{256, "AEAA"},
		{65535, "777Q"},
	}
	for _, tt := range tests {
		got, err := encodeID(tt.id)
		if err != nil {
			t.Fatalf("encodeID(%d): %v", tt.id, err)
		}
		if got != tt.want {
			t.Errorf("encodeID(%d) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

func TestDecodeKnownValues(t *testing.T) {
	tests := []struct {
		s    string
		want int64
	}{
		{"AAAA", 0},
		{"AAAQ", 1},
		{"AEAA", 256},
		{"777Q", 65535},
	}
	for _, tt := range tests {
		got, err := decodeID(tt.s)
		if err != nil {
			t.Fatalf("decodeID(%q): %v", tt.s, err)
		}
		if got != tt.want {
			t.Errorf("decodeID(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestDecodeCaseInsensitive(t *testing.T) {
	got, err := decodeID("aaaq")
	if err != nil {
		t.Fatalf("decodeID(aaaq): %v", err)
	}
	if got != 1 {
		t.Errorf("decodeID(aaaq) = %d, want 1", got)
	}
}

func TestRoundTrip(t *testing.T) {
	for _, id := range []int64{0, 1, 100, 256, 1000, 65535, 100000, 1000000} {
		s, err := encodeID(id)
		if err != nil {
			t.Fatalf("encodeID(%d): %v", id, err)
		}
		got, err := decodeID(s)
		if err != nil {
			t.Fatalf("decodeID(%q): %v", s, err)
		}
		if got != id {
			t.Errorf("roundtrip(%d): got %d", id, got)
		}
	}
}

func TestEncodeOutOfRange(t *testing.T) {
	if _, err := encodeID(-1); err == nil {
		t.Error("encodeID(-1) should fail")
	}
}

func TestDecodeInvalid(t *testing.T) {
	invalids := []string{"", "A", "!!!!", "AAAAAA"}
	for _, s := range invalids {
		if _, err := decodeID(s); err == nil {
			t.Errorf("decodeID(%q) should fail", s)
		}
	}
}
