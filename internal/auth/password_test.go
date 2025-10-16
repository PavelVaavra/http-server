package auth

import (
	"testing"
)

func TestPasswordHash(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "04234",
			expected: true,
		},
		// add more cases here
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("HashPassword(\"%v\") returns an error %v", c.input, err.Error())
		}
		match, err := CheckPasswordHash(c.input, hash)
		if err != nil {
			t.Errorf("CheckPasswordHash(\"%v\", \"%v\") returns an error %v", c.input, hash, err.Error())
		}
		if match != c.expected {
			t.Errorf("%v != %v", match, c.expected)
		}
	}
}
