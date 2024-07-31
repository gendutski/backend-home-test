package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertTo24HoursTime(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		for valid, source := range map[string]string{
			"18:30:00": "06:30:00 PM",
			"00:10:04": "12:10:04 AM",
			"08:27:30": "08:27:30 AM",
		} {
			result, err := ConvertTo24HoursTime(source)
			assert.Nil(t, err)
			assert.Equal(t, valid, result)
		}
	})
}
