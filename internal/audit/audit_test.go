package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRotateReleasesOldHandle reproduces bug9: Rotate overwrote l.f without
// closing the previous handle, so the old audit path stayed open for the
// process lifetime. On the CI builder the unlink of the rotated path failed
// ("file in use"), rotated garbage piled up on the artifacts disk, and the
// pipeline OOMed. After the fix the old path must be deletable and new logs
// must land in the rotated path only.
func TestRotateReleasesOldHandle(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old", "audit.log")
	newPath := filepath.Join(dir, "new", "audit.log")

	l, err := Open(oldPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := l.Log("put", "first"); err != nil {
		t.Fatalf("Log old: %v", err)
	}

	if err := l.Rotate(newPath); err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if got := l.Path(); got != newPath {
		t.Fatalf("Path after rotate = %q, want %q", got, newPath)
	}

	// Writes must go to the new path, never the old one.
	before, err := os.ReadFile(oldPath)
	if err != nil {
		t.Fatalf("read old: %v", err)
	}
	if err := l.Log("put", "second"); err != nil {
		t.Fatalf("Log new: %v", err)
	}
	after, err := os.ReadFile(oldPath)
	if err != nil {
		t.Fatalf("read old after: %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("old audit path changed after rotate: was %q now %q", before, after)
	}
	gotNew, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("read new: %v", err)
	}
	if !strings.Contains(string(gotNew), "second") {
		t.Fatalf("new audit path missing post-rotate entry, got %q", gotNew)
	}

	// The rotated-out path must be removable: this is the unlink the CI
	// cleanup task runs and that used to fail because the handle leaked.
	if err := os.Remove(oldPath); err != nil {
		t.Fatalf("Remove old path after rotate: %v", err)
	}

	if err := l.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestRotateClosesPreviousHandle ensures Rotate is safe to call repeatedly:
// every prior handle must be released, not just the original Open() one.
func TestRotateClosesPreviousHandle(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.log")
	second := filepath.Join(dir, "b.log")
	third := filepath.Join(dir, "c.log")

	l, err := Open(first)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := l.Rotate(second); err != nil {
		t.Fatalf("Rotate second: %v", err)
	}
	if err := l.Rotate(third); err != nil {
		t.Fatalf("Rotate third: %v", err)
	}
	// All intermediate paths must be deletable.
	for _, p := range []string{first, second} {
		if err := os.Remove(p); err != nil {
			t.Fatalf("Remove %s after repeat rotate: %v", p, err)
		}
	}
	_ = l.Close()
}
