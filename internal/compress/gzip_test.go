package compress

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressSuccess(t *testing.T) {
	in := []byte(`test`)

	cIn, err := Compress(in)
	require.NoError(t, err)
	assert.NotEmpty(t, cIn)
}

func TestDecompressSuccess(t *testing.T) {
	in := []byte(`test`)

	cIn, err := Compress(in)
	require.NoError(t, err)
	assert.NotEmpty(t, cIn)

	out, err := GzipDecompress(cIn)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}
