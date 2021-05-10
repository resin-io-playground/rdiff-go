package librsync

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignature(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	for _, tt := range allTestCases {
		t.Run(tt, func(t *testing.T) {
			file, magic, blockLen, strongLen, err := argsFromTestName(tt)
			r.NoError(err)

			inputData, err := ioutil.ReadFile("testdata/" + file + ".old")
			r.NoError(err)
			input := bytes.NewReader(inputData)

			output := &bytes.Buffer{}
			gotSig, err := Signature(input, output, blockLen, strongLen, magic)
			r.NoError(err)

			wantSig, err := readSignatureFile("testdata/" + tt + ".signature")
			r.NoError(err)
			a.Equal(wantSig.blockLen, gotSig.blockLen)
			a.Equal(wantSig.sigType, gotSig.sigType)
			a.Equal(wantSig.strongLen, gotSig.strongLen)

			outputData, err := ioutil.ReadAll(output)
			r.NoError(err)
			expectedData, err := ioutil.ReadFile("testdata/" + tt + ".signature")
			r.NoError(err)
			a.Equal(expectedData, outputData)
		})
	}
}
