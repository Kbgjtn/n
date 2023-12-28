package util

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhereAmI(t *testing.T) {
	testCwd := WhereAmI()
	slog.Info(testCwd)
	location := WhereAmI()
	assert.NotNil(t, location)
	assert.NotEmpty(t, location)
	assert.Equal(t, location, testCwd)
}
