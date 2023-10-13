package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const targetURL = `https://practicum.yandex.ru/`

func TestAddUrlAndGetUrlSuccess(t *testing.T) {
	r := new(ShortName).Init()

	err := r.AddUrl(`shortName`, targetURL)
	require.NoError(t, err)

	url, err := r.GetUrl(`shortName`)
	require.NoError(t, err)

	assert.Equal(t, targetURL, url)
}

func TestAddUrlShortNameExistsError(t *testing.T) {
	r := new(ShortName).Init()

	err := r.AddUrl(`shortName`, targetURL)
	require.NoError(t, err)

	err = r.AddUrl(`shortName`, targetURL)
	require.Error(t, err)
}

func TestAddUrlUrlExistsError(t *testing.T) {
	r := new(ShortName).Init()

	err := r.AddUrl(`shortName`, targetURL)
	require.NoError(t, err)

	err = r.AddUrl(`otherShortName`, targetURL)
	require.Error(t, err)
}

func TestGetUrlError(t *testing.T) {
	r := new(ShortName).Init()

	_, err := r.GetUrl(`shortName`)
	require.Error(t, err)
}
