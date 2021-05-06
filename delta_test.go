package librsync

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelta(t *testing.T) {
	tests := []struct {
		name      string
		blockLen  uint32
		strongLen uint32
		sigType   MagicNumber
		oldFile   string
		newFile   string
		wantDelta string
	}{
		{
			name:      "000 / Blake2 / 512 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  512,
			strongLen: 32,
			oldFile:   "testdata/delta/000.old",
			newFile:   "testdata/delta/000.new",
			wantDelta: "testdata/delta/000-blake2-512-32.delta",
		},
		{
			name:      "001 / Blake2 / 512 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  512,
			strongLen: 32,
			oldFile:   "testdata/delta/001.old",
			newFile:   "testdata/delta/001.new",
			wantDelta: "testdata/delta/001-blake2-512-32.delta",
		},
		{
			name:      "002 / Blake2 / 512 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  512,
			strongLen: 32,
			oldFile:   "testdata/delta/002.old",
			newFile:   "testdata/delta/002.new",
			wantDelta: "testdata/delta/002-blake2-512-32.delta",
		},
		{
			name:      "003 / Blake2 / 512 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  512,
			strongLen: 32,
			oldFile:   "testdata/delta/003.old",
			newFile:   "testdata/delta/003.new",
			wantDelta: "testdata/delta/003-blake2-512-32.delta",
		},
		{
			name:      "004 / Blake2 / 512 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  512,
			strongLen: 32,
			oldFile:   "testdata/delta/004.old",
			newFile:   "testdata/delta/004.new",
			wantDelta: "testdata/delta/004-blake2-512-32.delta",
		},
		{
			name:      "004 / Blake2 / 2222 / 31",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  2222,
			strongLen: 31,
			oldFile:   "testdata/delta/004.old",
			newFile:   "testdata/delta/004.new",
			wantDelta: "testdata/delta/004-blake2-2222-31.delta",
		},
		{
			name:      "004 / Blake2 / 1024 / 28",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  1024,
			strongLen: 28,
			oldFile:   "testdata/delta/004.old",
			newFile:   "testdata/delta/004.new",
			wantDelta: "testdata/delta/004-blake2-1024-28.delta",
		},
		{
			name:      "005 / Blake2 / 512 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  512,
			strongLen: 32,
			oldFile:   "testdata/delta/005.old",
			newFile:   "testdata/delta/005.new",
			wantDelta: "testdata/delta/005-blake2-512-32.delta",
		},
		{
			name:      "006 / Blake2 / 2 / 32",
			sigType:   BLAKE2_SIG_MAGIC,
			blockLen:  2,
			strongLen: 32,
			oldFile:   "testdata/delta/006.old",
			newFile:   "testdata/delta/006.new",
			wantDelta: "testdata/delta/006-blake2-2-32.delta",
		},
	}

	r := require.New(t)
	a := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFile, err := os.Open(tt.oldFile)
			r.NoError(err)

			sig, err := Signature(oldFile, io.Discard, tt.blockLen, tt.strongLen, tt.sigType)
			r.NoError(err)

			newFile, err := os.Open(tt.newFile)
			r.NoError(err)
			output := &bytes.Buffer{}

			err = Delta(sig, newFile, output)
			r.NoError(err)

			outputDelta, err := ioutil.ReadAll(output)
			r.NoError(err)
			expectedDelta, err := ioutil.ReadFile(tt.wantDelta)
			r.NoError(err)
			a.Equal(expectedDelta, outputDelta)
		})
	}
}
