package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CountIntegerInFile(t *testing.T) {
	t.Run("testfile.txt", func(t *testing.T) {
		totalNumbers, sumNumbers, err := CountIntegerInFile("testfile.txt")

		assert.Nil(t, err)
		assert.Equal(t, 5, totalNumbers)
		assert.Equal(t, 15, sumNumbers)
	})

	t.Run("testfile2.txt", func(t *testing.T) {
		totalNumbers, sumNumbers, err := CountIntegerInFile("testfile2.txt")

		assert.Nil(t, err)
		assert.Equal(t, 10, totalNumbers)
		assert.Equal(t, 55, sumNumbers)
	})
}
