package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSHA256(t *testing.T) {
	hash, err := md5sum("testdata/dummy.txt")
	require.NoError(t, err)
	require.Equal(t, "06ad47d8e64bd28de537b62ff85357c4", hash)
}
