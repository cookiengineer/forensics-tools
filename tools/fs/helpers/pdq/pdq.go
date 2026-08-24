package pdq

import (
	"bytes"
	"encoding/hex"
	"image"
	"math/bits"
	"os"

	// Register image decoders so image.Decode can handle every format the
	// fs-cleanup tool classifies as an image.
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
)

// hasher is the shared PDQ (Facebook/Meta perceptual image hash) instance. Its
// Calculate method has a value receiver and holds no mutable state, so a single
// instance is safe for concurrent use.
var hasher = mustNew()

func mustNew() imghash.PDQ {
	h, err := imghash.NewPDQ()
	if err != nil {
		panic(err)
	}
	return h
}

// HashFile computes the 256-bit PDQ hash of an image file and returns its
// decoded pixel dimensions. The hash is encoded as a lowercase hex string.
func HashFile(path string) (hash string, width, height int, err error) {
	img, _, err := decodeFile(path)
	if err != nil {
		return "", 0, 0, err
	}
	hash, err = hashImage(img)
	if err != nil {
		return "", 0, 0, err
	}
	b := img.Bounds()
	return hash, b.Dx(), b.Dy(), nil
}

// HashBytes computes the 256-bit PDQ hash from the encoded bytes of an image
// (e.g. a frame extracted by ffmpeg), returned as a lowercase hex string.
func HashBytes(data []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return hashImage(img)
}

func decodeFile(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func hashImage(img image.Image) (string, error) {
	h, err := hasher.Calculate(img)
	if err != nil {
		return "", err
	}
	return toHex(h), nil
}

func toHex(h hashtype.Hash) string {
	return hex.EncodeToString([]byte(h.(hashtype.Binary)))
}

// Similarity returns the similarity of two hex-encoded PDQ hashes as a
// percentage in the range [0, 100], derived from their Hamming distance across
// the 256 hash bits. 100 means identical; 0 means every bit differs.
func Similarity(a, b string) float64 {
	ha, errA := hex.DecodeString(a)
	hb, errB := hex.DecodeString(b)
	if errA != nil || errB != nil || len(ha) != len(hb) {
		return 0
	}
	total := len(ha) * 8
	if total == 0 {
		return 100
	}
	dist := 0
	for i := range ha {
		dist += bits.OnesCount8(ha[i] ^ hb[i])
	}
	return (1 - float64(dist)/float64(total)) * 100
}
