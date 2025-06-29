package routes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPathParam(t *testing.T) {
	require.True(t, isPathParam("{id}"))
	require.False(t, isPathParam("*id"))
}

func TestIsWildcard(t *testing.T) {
	require.True(t, isWildcard("*id"))
	require.False(t, isWildcard("{id}"))
}
