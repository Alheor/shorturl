package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const targetURL = `https://practicum.yandex.ru/`

func TestAddURLAndGetURLSuccess(t *testing.T) {
	r := new(ShortName).Init()

	err := r.AddURL(`shortName`, targetURL)
	require.NoError(t, err)

	url, err := r.GetURL(`shortName`)
	require.NoError(t, err)

	assert.Equal(t, targetURL, url)
}

func TestAddURLShortNameExistsError(t *testing.T) {
	r := new(ShortName).Init()

	err := r.AddURL(`shortName`, targetURL)
	require.NoError(t, err)

	err = r.AddURL(`shortName`, targetURL)
	require.Error(t, err)
}

func TestAddURLURLExistsError(t *testing.T) {
	r := new(ShortName).Init()

	err := r.AddURL(`shortName`, targetURL)
	require.NoError(t, err)

	err = r.AddURL(`otherShortName`, targetURL)
	require.Error(t, err)
}

func TestGetURLError(t *testing.T) {
	r := new(ShortName).Init()

	_, err := r.GetURL(`shortName`)
	require.Error(t, err)
}
