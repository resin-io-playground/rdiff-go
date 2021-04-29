package librsync

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func signature(b *testing.B, src io.Reader) *SignatureType {
	var (
		magic            = BLAKE2_SIG_MAGIC
		blockLen  uint32 = 512
		strongLen uint32 = 32
		bufSize          = 65536
	)

	s, err := Signature(
		bufio.NewReaderSize(src, bufSize),
		ioutil.Discard,
		blockLen, strongLen, magic)
	if err != nil {
		b.Error(err)
	}

	return s
}

func benchmarkSignature(b *testing.B, totalBytes int64) {
	b.SetBytes(totalBytes)

	for i := 0; i < b.N; i++ {
		src := io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes)
		signature(b, src)
	}
}

func BenchmarkSignature1MB(b *testing.B) {
	benchmarkSignature(b, 1_000_000)
}

func BenchmarkSignature1GB(b *testing.B) {
	benchmarkSignature(b, 1_000_000_000)
}

func benchmarkDeltaChangedTail(b *testing.B, totalBytes int64) {
	newBytes := totalBytes / 10
	oldBytes := totalBytes - newBytes
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)
	s := signature(b, oldData)

	b.SetBytes(totalBytes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		newSeed := time.Now().UnixNano()
		newData := io.MultiReader(
			io.LimitReader(rand.New(rand.NewSource(oldSeed)), oldBytes),
			io.LimitReader(rand.New(rand.NewSource(newSeed)), newBytes),
		)

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaChangedTail1GB(b *testing.B) {
	benchmarkDeltaChangedTail(b, 1_000_000_000)
}

func BenchmarkDeltaChangedTail1MB(b *testing.B) {
	benchmarkDeltaChangedTail(b, 1_000_000)
}
