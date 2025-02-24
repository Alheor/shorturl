package urlhasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {

	randomShortName := new(ShortName)

	val1 := randomShortName.Generate()
	val2 := randomShortName.Generate()

	assert.NotEmpty(t, val1)
	assert.NotEmpty(t, val2)
	assert.NotEqual(t, val1, val2)
}
