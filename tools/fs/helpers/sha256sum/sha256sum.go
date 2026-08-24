package sha256sum

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Sum returns the lowercase hex SHA-256 digest of the given bytes using the
// standard library. This is used for the tiny normalized frame buffers produced
// by ffmpeg, where spawning an external process would be wasteful.
func Sum(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// hashStdlib streams the file through the standard-library SHA-256
// implementation, keeping memory usage constant regardless of file size.
func hashStdlib(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	buf := make([]byte, 1<<20)
	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashFile computes the SHA-256 of the file at path. It shells out to the
// system `sha256sum` binary (usually present on Linux and typically faster than
// a Go implementation for large files). If the binary is unavailable or fails,
// it falls back to the standard library.
func HashFile(path string) (string, error) {
	bin, err := exec.LookPath("sha256sum")
	if err != nil {
		return hashStdlib(path)
	}

	out, err := exec.Command(bin, path).Output()
	if err != nil {
		return hashStdlib(path)
	}

	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return hashStdlib(path)
	}
	return strings.ToLower(fields[0]), nil
}
