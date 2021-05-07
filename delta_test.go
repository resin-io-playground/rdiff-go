package librsync

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelta(t *testing.T) {
	tests := []string{
		"000-blake2-11-23",
		"000-blake2-512-32",
		"000-md4-256-7",
		"001-blake2-512-32",
		"001-blake2-776-31",
		"001-md4-777-15",
		"002-blake2-512-32",
		"002-blake2-431-19",
		"002-md4-128-16",
		"003-blake2-512-32",
		"003-blake2-1024-13",
		"003-md4-1024-13",
		"004-blake2-1024-28",
		"004-blake2-2222-31",
		"004-blake2-512-32",
		"005-blake2-512-32",
		"005-blake2-1000-18",
		"005-md4-999-14",
		"006-blake2-2-32",
		"007-blake2-5-32",
		"007-blake2-4-32",
		"007-blake2-3-32",
	}

	r := require.New(t)
	a := assert.New(t)

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			file, _, _, _, err := argsFromTestName(tt)
			r.NoError(err)

			sig, err := readSignatureFile("testdata/" + tt + ".signature")
			r.NoError(err)

			newFile, err := os.Open("testdata/" + file + ".new")
			r.NoError(err)
			output := &bytes.Buffer{}

			err = Delta(sig, newFile, output)
			r.NoError(err)

			gotDelta, err := ioutil.ReadAll(output)
			r.NoError(err)
			wantDelta, err := ioutil.ReadFile("testdata/" + tt + ".delta")
			r.NoError(err)
			a.Equal(wantDelta, gotDelta)
		})
	}
}
