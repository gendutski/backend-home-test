package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsPalindrom(t *testing.T) {
	t.Run("palindrom words", func(t *testing.T) {
		assert.True(t, IsPalindrom("A man, a plan, a canal: Panama "))
		assert.True(t, IsPalindrom("Was it a car or a cat I saw?"))
		assert.True(t, IsPalindrom("Mr. Owl ate my metal worm"))
		assert.True(t, IsPalindrom("Doc, note: I dissent. A fast never prevents a fatness. I diet on cod."))
	})

	t.Run("non palindrom words", func(t *testing.T) {
		assert.False(t, IsPalindrom("Malam ini aku terbangun"))
		assert.False(t, IsPalindrom("Buatlah sebuah program dengan Golang"))
		assert.False(t, IsPalindrom("Ibu ratna antar ibu"))
	})
}
