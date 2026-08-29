package sandbox

import "testing"

func TestValidateID(t *testing.T) {
	t.Parallel()

	for _, id := range []string{"baseline-1", "run_002", "a.b"} {
		if err := ValidateID(id); err != nil {
			t.Fatalf("ValidateID(%q) returned %v", id, err)
		}
	}

	for _, id := range []string{"", "../escape", "has space", "a/b"} {
		if err := ValidateID(id); err == nil {
			t.Fatalf("ValidateID(%q) unexpectedly succeeded", id)
		}
	}
}
