package token

import (
	"slices"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("returns a valid token", func(t *testing.T) {
		// when
		result := New()
		// then
		assert.True(t, Valid(result), "invalid token %q", result)
	})

	t.Run("encodes a uuid v4", func(t *testing.T) {
		// given
		token := New()
		// when
		result, err := ToUUID(token)
		// then
		assert.NoError(t, err)
		assert.Equal(t, byte(4), result[6]>>4, "version of %s", result)
		assert.Equal(t, byte(0b10), result[8]>>6, "variant of %s", result)
	})

	t.Run("returns a different token on each call", func(t *testing.T) {
		// when
		result1 := New()
		result2 := New()
		// then
		assert.NotEqual(t, result1, result2)
	})
}

func TestFromUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    uuid.UUID
		expected string
	}{
		{
			name:     "nil uuid encodes as all zeros",
			input:    uuid.Nil(),
			expected: "0000000000000000000000000",
		},
		{
			name:     "max uuid encodes as the largest token",
			input:    uuid.Max(),
			expected: "f5lxx1zz5pnorynqglhzmsp33",
		},
		{
			name:     "value one is left-padded with zeros",
			input:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			expected: "0000000000000000000000001",
		},
		{
			name:     "value 36 carries into the second digit",
			input:    uuid.MustParse("00000000-0000-0000-0000-000000000024"),
			expected: "0000000000000000000000010",
		},
		{
			name:     "value 2^64 carries from the low half into the high half",
			input:    uuid.MustParse("00000000-0000-0001-0000-000000000000"),
			expected: "0000000000003w5e11264sgsg",
		},
		{
			name:     "matches the lowercase narciso token for the same uuid",
			input:    uuid.MustParse("41c9ad60-0cab-4e1b-afc8-cf97fbb94662"),
			expected: "3w7nni025418b96lydzqxyb8i",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result := FromUUID(input)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}

	t.Run("tokens sort in the same order as their uuids", func(t *testing.T) {
		// given
		input := []uuid.UUID{
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			uuid.MustParse("00000000-0000-0000-0000-00000000000a"),
			uuid.MustParse("00000000-0000-0000-ffff-ffffffffffff"),
			uuid.MustParse("00000000-0000-0001-0000-000000000000"),
			uuid.MustParse("0197a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b"),
			uuid.MustParse("0197a1b2-c3d5-7000-8000-000000000000"),
			uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		}
		// when
		result := make([]string, len(input))
		for i, u := range input {
			result[i] = FromUUID(u)
		}
		// then
		assert.True(t, slices.IsSortedFunc(input, uuid.UUID.Compare), "fixture must be sorted")
		assert.True(t, slices.IsSorted(result), "tokens out of order: %s", strings.Join(result, " "))
	})
}

func TestToUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uuid.UUID
	}{
		{
			name:     "all zeros decodes as nil uuid",
			input:    "0000000000000000000000000",
			expected: uuid.Nil(),
		},
		{
			name:     "largest token decodes as max uuid",
			input:    "f5lxx1zz5pnorynqglhzmsp33",
			expected: uuid.Max(),
		},
		{
			name:     "last digit one decodes as value one",
			input:    "0000000000000000000000001",
			expected: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
		{
			name:     "second digit one decodes as value 36",
			input:    "0000000000000000000000010",
			expected: uuid.MustParse("00000000-0000-0000-0000-000000000024"),
		},
		{
			name:     "value 2^64 carries from the low half into the high half",
			input:    "0000000000003w5e11264sgsg",
			expected: uuid.MustParse("00000000-0000-0001-0000-000000000000"),
		},
		{
			name:     "matches the uuid of the lowercase narciso token",
			input:    "3w7nni025418b96lydzqxyb8i",
			expected: uuid.MustParse("41c9ad60-0cab-4e1b-afc8-cf97fbb94662"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result, err := ToUUID(input)
			// then
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}

	invalid := []struct {
		name  string
		input string
	}{
		{name: "one above the largest 128-bit value", input: "f5lxx1zz5pnorynqglhzmsp34"},
		{name: "uppercase letter", input: "3W7nni025418b96lydzqxyb8i"},
		{name: "empty string", input: ""},
		{name: "one character short", input: "000000000000000000000000"},
		{name: "one character long", input: "00000000000000000000000000"},
		{name: "hyphenated uuid", input: "41c9ad60-0cab-4e1b-afc8-cf97fbb94662"},
		{name: "character outside the alphabet", input: "000000000000000000000000_"},
	}

	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result, err := ToUUID(input)
			// then
			assert.ErrorIs(t, err, ErrInvalidToken)
			assert.Equal(t, uuid.Nil(), result)
		})
	}

	t.Run("reverses FromUUID for every uuid", func(t *testing.T) {
		// given
		input := []uuid.UUID{
			uuid.Nil(),
			uuid.MustParse("00000000-0000-0000-ffff-ffffffffffff"),
			uuid.MustParse("0197a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b"),
			uuid.MustParse("41c9ad60-0cab-4e1b-afc8-cf97fbb94662"),
			uuid.Max(),
		}
		// when
		result := make([]uuid.UUID, len(input))
		for i, u := range input {
			result[i], _ = ToUUID(FromUUID(u))
		}
		// then
		assert.Equal(t, input, result)
	})
}

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "all zeros", input: "0000000000000000000000000", expected: true},
		{name: "largest 128-bit value", input: "f5lxx1zz5pnorynqglhzmsp33", expected: true},
		{name: "digits and lowercase letters", input: "3w7nni025418b96lydzqxyb8i", expected: true},
		{name: "one above the largest 128-bit value", input: "f5lxx1zz5pnorynqglhzmsp34", expected: false},
		{name: "above range in the first character", input: "g000000000000000000000000", expected: false},
		{name: "all z", input: "zzzzzzzzzzzzzzzzzzzzzzzzz", expected: false},
		{name: "empty string", input: "", expected: false},
		{name: "one character short", input: "000000000000000000000000", expected: false},
		{name: "one character long", input: "00000000000000000000000000", expected: false},
		{name: "uppercase letter", input: "3W7nni025418b96lydzqxyb8i", expected: false},
		{name: "colon just after 9 in ASCII", input: "000000000000000000000000:", expected: false},
		{name: "backtick just before a in ASCII", input: "000000000000000000000000`", expected: false},
		{name: "slash just before 0 in ASCII", input: "000000000000000000000000/", expected: false},
		{name: "brace just after z in ASCII", input: "000000000000000000000000{", expected: false},
		{name: "hyphenated uuid", input: "41c9ad60-0cab-4e1b-afc8-cf97fbb94662", expected: false},
		{name: "multi-byte character with 25 bytes", input: "00000000000000000000000é", expected: false},
		{name: "embedded NUL byte", input: "000000000000\x00000000000000", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result := Valid(input)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}

	t.Run("accepts every token produced from a uuid", func(t *testing.T) {
		// given
		input := FromUUID(uuid.MustParse("0197a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b"))
		// when
		result := Valid(input)
		// then
		assert.True(t, result)
	})
}

func TestAllocations(t *testing.T) {
	u := uuid.MustParse("41c9ad60-0cab-4e1b-afc8-cf97fbb94662")

	tests := []struct {
		name     string
		call     func()
		expected float64
	}{
		{name: "New allocates only the result string", call: func() { _ = New() }, expected: 1},
		{name: "FromUUID allocates only the result string", call: func() { _ = FromUUID(u) }, expected: 1},
		{name: "ToUUID does not allocate", call: func() { _, _ = ToUUID("3w7nni025418b96lydzqxyb8i") }, expected: 0},
		{name: "ToUUID with an invalid token does not allocate", call: func() { _, _ = ToUUID("not a token") }, expected: 0},
		{name: "Valid does not allocate", call: func() { _ = Valid("3w7nni025418b96lydzqxyb8i") }, expected: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			call := tc.call
			// when
			result := testing.AllocsPerRun(100, call)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}
