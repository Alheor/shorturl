package randomname

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSuccess(t *testing.T) {

	randomShortName := new(ShortName)

	val1 := randomShortName.Generate()
	val2 := randomShortName.Generate()

	assert.NotEqual(t, val1, val2)
}
