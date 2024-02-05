package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMajorVersion(t *testing.T) {
	v, err := getMajorVersion("1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	v, err = getMajorVersion("0.0.0")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("65535.0.0")
	assert.NoError(t, err)
	assert.Equal(t, int64(65535), v)
}

func TestGetMajorVersionError(t *testing.T) {
	v, err := getMajorVersion("a")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("1.1.1.1")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("-1.0.0")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("184467440737095516150.0.0")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)
}
