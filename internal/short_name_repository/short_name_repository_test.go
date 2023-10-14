package short_name_repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const targetURL = `https://practicum.yandex.ru/`

func TestAddURLAndGetURLSuccess(t *testing.T) {
	r := Init()

	err := r.Add(`shortName`, targetURL)
	require.NoError(t, err)

	url, err := r.Get(`shortName`)
	require.NoError(t, err)

	assert.Equal(t, targetURL, url)
}

func TestAddURLShortNameExistsError(t *testing.T) {
	r := Init()

	err := r.Add(`shortName`, targetURL)
	require.NoError(t, err)

	err = r.Add(`shortName`, targetURL)
	require.Error(t, err)
}

func TestAddURLURLExistsError(t *testing.T) {
	r := Init()

	err := r.Add(`shortName`, targetURL)
	require.NoError(t, err)

	err = r.Add(`otherShortName`, targetURL)
	require.Error(t, err)
}

func TestGetURLError(t *testing.T) {
	r := Init()

	_, err := r.Get(`shortName`)
	require.Error(t, err)
}
