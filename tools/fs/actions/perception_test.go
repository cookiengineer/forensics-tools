package actions

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"fs/helpers/ffmpeg"
)

func runCmd(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}

func TestPerceptionHashImagesIntegration(t *testing.T) {
	if !ffmpeg.Available() {
		t.Skip("ffmpeg/ffprobe not available")
	}
	dir := t.TempDir()
	base := filepath.Join(dir, "base.png")
	runCmd(t, dir, "ffmpeg", "-v", "error", "-f", "lavfi", "-i", "color=c=red:s=64x64", "-frames:v", "1", base)

	// Same pixels, different bytes (extra trailing data is ignored by decoders),
	// sizes within 99% of each other.
	raw, err := os.ReadFile(base)
	if err != nil {
		t.Fatal(err)
	}
	variant := filepath.Join(dir, "a.png")
	if err := os.WriteFile(variant, append(raw, []byte("X")...), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	Hash(files)
	PerceptionHash(files)

	baseFile := files["base.png"]
	variantFile := files["a.png"]
	if baseFile.Delete {
		t.Error("base.png should be preserved")
	}
	if !variantFile.Delete {
		t.Fatal("a.png should be marked for deletion")
	}
	if variantFile.DuplicateOf != "base.png" {
		t.Errorf("DuplicateOf = %q, want base.png", variantFile.DuplicateOf)
	}
	if !variantFile.Perceptual {
		t.Error("a.png should be flagged as a perceptual (not exact) duplicate")
	}
	if variantFile.Similarity != 100 {
		t.Errorf("Similarity = %v, want 100 (identical PDQ hash)", variantFile.Similarity)
	}
}

func TestPerceptionHashVideosIntegration(t *testing.T) {
	if !ffmpeg.Available() {
		t.Skip("ffmpeg/ffprobe not available")
	}
	dir := t.TempDir()
	base := filepath.Join(dir, "base.mp4")
	runCmd(t, dir, "ffmpeg", "-v", "error", "-f", "lavfi", "-i", "testsrc=duration=3:size=128x128:rate=10",
		"-pix_fmt", "yuv420p", base)
	// Same content, same duration, but different bytes via a metadata change.
	runCmd(t, dir, "ffmpeg", "-v", "error", "-i", base, "-c", "copy", "-metadata", "title=hello", "a.mp4")

	files, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	Hash(files)
	PerceptionHash(files)

	if files["base.mp4"].Delete {
		t.Error("base.mp4 should be preserved")
	}
	variant := files["a.mp4"]
	if !variant.Delete {
		t.Fatal("a.mp4 should be marked for deletion")
	}
	if variant.DuplicateOf != "base.mp4" {
		t.Errorf("DuplicateOf = %q, want base.mp4", variant.DuplicateOf)
	}
}

// TestPerceptionHashDistinctVideosIntegration uses the real fixtures in the
// module testdata directory: two unrelated videos with similar file sizes but
// different durations and content, plus an exact byte-for-byte copy of one of
// them. The two distinct videos must not be reported as near-identical.
func TestPerceptionHashDistinctVideosIntegration(t *testing.T) {
	if !ffmpeg.Available() {
		t.Skip("ffmpeg/ffprobe not available")
	}
	dir := filepath.Join("..", "testdata")
	if _, err := os.Stat(filepath.Join(dir, "video-original.mp4")); err != nil {
		t.Skip("testdata videos not present")
	}

	files, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	Hash(files)
	PerceptionHash(files)

	distinct := files["video-within-same-byterange.mp4"]
	if distinct == nil {
		t.Fatal("video-within-same-byterange.mp4 not found")
	}
	if distinct.Delete {
		t.Errorf("distinct video %s was marked for deletion", distinct.Name)
	}
	if distinct.DuplicateOf == "video-original.mp4" {
		t.Errorf("distinct video %s was reported as a duplicate of video-original.mp4", distinct.Name)
	}

	// The byte-for-byte copy is a genuine exact duplicate and must be deleted.
	if copyFile := files["video-copy.mp4"]; copyFile == nil || !copyFile.Delete {
		t.Error("video-copy.mp4 (exact copy) should be marked for deletion")
	}
}

// TestPerceptionHashImagesResizeIntegration uses the real fixtures in testdata
// to confirm that a resized variant of an image (far smaller byte size, lower
// resolution) is still matched by its PDQ perceptual hash and marked for
// deletion in favour of the higher-resolution original.
func TestPerceptionHashImagesResizeIntegration(t *testing.T) {
	dir := filepath.Join("..", "testdata")
	if _, err := os.Stat(filepath.Join(dir, "image-smaller.jpeg")); err != nil {
		t.Skip("testdata images not present")
	}

	files, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	Hash(files)
	PerceptionHash(files)

	smaller := files["image-smaller.jpeg"]
	if smaller == nil {
		t.Fatal("image-smaller.jpeg not found")
	}
	if !smaller.Delete {
		t.Fatal("image-smaller.jpeg (resized) should be marked for deletion")
	}
	if smaller.DuplicateOf != "image-original.jpeg" {
		t.Errorf("DuplicateOf = %q, want image-original.jpeg", smaller.DuplicateOf)
	}
	if !smaller.Perceptual || smaller.Similarity != 100 {
		t.Errorf("smaller should be a perceptual duplicate with 100%% PDQ similarity, got Perceptual=%v Similarity=%v",
			smaller.Perceptual, smaller.Similarity)
	}
	if original := files["image-original.jpeg"]; original == nil || original.Delete {
		t.Error("image-original.jpeg (higher resolution) should be preserved")
	}
}
