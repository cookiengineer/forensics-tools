package actions

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExtensionAndClassify(t *testing.T) {
	cases := map[string]struct {
		ext   string
		media MediaKind
	}{
		"photo.JPG": {"jpg", MediaImage},
		"pic.png":   {"png", MediaImage},
		"movie.MP4": {"mp4", MediaVideo},
		"clip.webm": {"webm", MediaVideo},
		"doc.txt":   {"txt", MediaNone},
		"noext":     {"", MediaNone},
	}
	for name, want := range cases {
		if got := Extension(name); got != want.ext {
			t.Errorf("Extension(%q) = %q, want %q", name, got, want.ext)
		}
		if got := ClassifyMedia(want.ext); got != want.media {
			t.Errorf("ClassifyMedia(%q) = %v, want %v", want.ext, got, want.media)
		}
	}
}

func TestSimilarityPercent(t *testing.T) {
	cases := []struct {
		a, b int64
		want float64
	}{
		{100, 100, 100},
		{99, 100, 99},
		{100, 99, 99},
		{0, 0, 100},
		{0, 100, 0},
	}
	for _, c := range cases {
		if got := SimilarityPercent(c.a, c.b); got != c.want {
			t.Errorf("SimilarityPercent(%d,%d) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestSizeWithinPercent(t *testing.T) {
	if !SizeWithinPercent(99, 100, 99) {
		t.Error("99 vs 100 should be within 99%")
	}
	if SizeWithinPercent(98, 100, 99) {
		t.Error("98 vs 100 should not be within 99%")
	}
	if !SizeWithinPercent(100, 100, 99) {
		t.Error("equal sizes should be within 99%")
	}
}

func TestSameSecond(t *testing.T) {
	cases := []struct {
		a, b float64
		want bool
	}{
		{3.0, 3.0, true},
		{3.1, 3.4, true},
		{3.0, 4.0, false},
		{2.9, 3.1, true},
	}
	for _, c := range cases {
		if got := SameSecond(c.a, c.b); got != c.want {
			t.Errorf("SameSecond(%v,%v) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestFrameSeeks(t *testing.T) {
	if got := startSeek(5.0); got != 1.0 {
		t.Errorf("startSeek(5) = %v, want 1.0", got)
	}
	if got := startSeek(1.0); got != 0.5 {
		t.Errorf("startSeek(1) = %v, want 0.5", got)
	}
	if seek, ok := endSeek(5.0); !ok || seek != 4.0 {
		t.Errorf("endSeek(5) = %v,%v, want 4.0,true", seek, ok)
	}
	if _, ok := endSeek(1.0); ok {
		t.Error("endSeek(1) should report ok=false")
	}
}

func TestOriginalityScore(t *testing.T) {
	original := originalityScore("report.txt")
	if s := originalityScore("report.txt.bak"); s <= original {
		t.Errorf("bak should score higher than original: %d vs %d", s, original)
	}
	if s := originalityScore("report-1.txt"); s <= original {
		t.Errorf("-1 should score higher than original: %d vs %d", s, original)
	}
	if s := originalityScore("report-new.txt"); s <= original {
		t.Errorf("-new should score higher than original: %d vs %d", s, original)
	}
	if originalityScore("file-a.txt") != originalityScore("file-g.txt") {
		t.Error("unrelated letter suffixes should tie")
	}
}

func TestPreserveOriginal(t *testing.T) {
	now := time.Now()
	mk := func(name string, mod time.Time) *File {
		return &File{Name: name, ModTime: mod}
	}

	// Filename identifies the original.
	kept := PreserveOriginal([]*File{
		mk("report.txt", now),
		mk("report-1.txt", now.Add(-time.Hour)), // even if older, it's a derivative
		mk("report.txt.bak", now.Add(-2*time.Hour)),
	})
	if kept.Name != "report.txt" {
		t.Errorf("preserved %q, want report.txt", kept.Name)
	}

	// Uninferable originality falls back to oldest modtime.
	oldest := now.Add(-time.Hour)
	kept = PreserveOriginal([]*File{
		mk("file-b.txt", now),
		mk("file-a.txt", oldest),
	})
	if kept.Name != "file-a.txt" {
		t.Errorf("preserved %q, want file-a.txt (oldest)", kept.Name)
	}

	// Equal modtime falls back to alphabetical.
	kept = PreserveOriginal([]*File{
		mk("z.txt", now),
		mk("a.txt", now),
	})
	if kept.Name != "a.txt" {
		t.Errorf("preserved %q, want a.txt (alphabetical)", kept.Name)
	}
}

func TestMarkShorter(t *testing.T) {
	a := &File{Name: "a.mp4", Size: 100, Duration: 5.0}
	b := &File{Name: "b.mp4", Size: 99, Duration: 3.0}
	markShorter(pair{a: a, b: b})
	if a.Delete {
		t.Error("longer video should be kept")
	}
	if !b.Delete || b.DuplicateOf != "a.mp4" {
		t.Errorf("shorter video should be deleted: %+v", b)
	}
}

func TestMarkNewerForDeletion(t *testing.T) {
	older := &File{Name: "old.mp4", ModTime: time.Now().Add(-time.Hour), Size: 100}
	newer := &File{Name: "new.mp4", ModTime: time.Now(), Size: 99}
	markNewerForDeletion(older, newer)
	if older.Delete {
		t.Error("older should be kept")
	}
	if !newer.Delete || newer.DuplicateOf != "old.mp4" {
		t.Errorf("newer should be deleted: %+v", newer)
	}
}

func TestScanAndHashIntegration(t *testing.T) {
	dir := t.TempDir()
	content := []byte("identical content\n")
	for _, name := range []string{"report.txt", "report-1.txt", "report.txt.bak"} {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "unrelated.txt"), []byte("totally different bytes here\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("expected 4 files, got %d", len(files))
	}

	Hash(files)

	deleted := 0
	kept := map[string]bool{}
	for _, f := range files {
		if f.Delete {
			deleted++
		} else {
			kept[f.Name] = true
		}
	}
	if deleted != 2 {
		t.Fatalf("expected 2 deleted, got %d", deleted)
	}
	if !kept["report.txt"] {
		t.Error("report.txt (the original) should be preserved")
	}
	if !kept["unrelated.txt"] {
		t.Error("unrelated.txt should be preserved")
	}
}
