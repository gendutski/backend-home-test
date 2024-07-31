package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateRandomSequence(t *testing.T) {
	validUsed := map[int]bool{
		1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true,
	}

	used := make(map[int]bool)
	seq := []int{}
	seq = CreateRandomSequence(seq, used)

	assert.Equal(t, 10, len(seq))
	assert.Equal(t, validUsed, used)
}
