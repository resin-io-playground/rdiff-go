package librsync

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignature(t *testing.T) {
	tests := []struct {
		name         string
		argBlockLen  uint32
		argStrongLen uint32
		argSigType   MagicNumber
		wantOutput   string
	}{
		{
			name:         "Blake2 / 512 / 32",
			argSigType:   BLAKE2_SIG_MAGIC,
			argBlockLen:  512,
			argStrongLen: 32,
			wantOutput:   "testdata/signature/blake2-512-32.sig",
		},
		{
			name:         "Blake2 / 2048 / 28",
			argSigType:   BLAKE2_SIG_MAGIC,
			argBlockLen:  2048,
			argStrongLen: 28,
			wantOutput:   "testdata/signature/blake2-2048-28.sig",
		},
		{
			name:         "Blake2 / 1171 / 31",
			argSigType:   BLAKE2_SIG_MAGIC,
			argBlockLen:  1171,
			argStrongLen: 31,
			wantOutput:   "testdata/signature/blake2-1171-31.sig",
		},
		{
			name:         "MD4 / 1111 / 15",
			argSigType:   MD4_SIG_MAGIC,
			argBlockLen:  1111,
			argStrongLen: 15,
			wantOutput:   "testdata/signature/md4-1111-15.sig",
		},
	}

	r := require.New(t)
	a := assert.New(t)

	inputData, err := ioutil.ReadFile("testdata/signature/signature.input")
	r.NoError(err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := bytes.NewReader(inputData)

			output := &bytes.Buffer{}
			got, err := Signature(input, output, tt.argBlockLen, tt.argStrongLen, tt.argSigType)
			r.NoError(err)

			a.Equal(tt.argBlockLen, got.blockLen)
			a.Equal(tt.argSigType, got.sigType)
			a.Equal(tt.argStrongLen, got.strongLen)

			outputData, err := ioutil.ReadAll(output)
			r.NoError(err)
			expectedData, err := ioutil.ReadFile(tt.wantOutput)
			r.NoError(err)
			a.Equal(expectedData, outputData)
		})
	}
}
