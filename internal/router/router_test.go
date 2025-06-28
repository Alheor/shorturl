package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSuccess(t *testing.T) {
	list := GetRoutes()

	assert.NotEmpty(t, list.Routes())
}
