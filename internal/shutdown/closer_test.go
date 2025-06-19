package shutdown

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSuccess(t *testing.T) {
	Init()
	assert.NotNil(t, closer)
}

func TestAddSuccess(t *testing.T) {
	Init()
	assert.NotNil(t, closer)

	f := func(ctx context.Context) error {
		return nil
	}

	closer.Add(f)

	assert.Len(t, closer.funcs, 1)
}

func TestCloseSuccess(t *testing.T) {
	Init()
	assert.NotNil(t, closer)

	val := 1

	f := func(ctx context.Context) error {
		val = 0
		return nil
	}

	closer.Add(f)
	assert.Equal(t, val, 1)

	closer.Close(context.Background())
	assert.Equal(t, val, 0)
}
