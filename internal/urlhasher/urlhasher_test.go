package urlhasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHashSuccess(t *testing.T) {
	in := `test`
	hash := GetHash(in)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, HashLength)
}
