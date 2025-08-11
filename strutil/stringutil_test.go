package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n\r", true},
		{"hello", false},
		{" hello ", false},
	}

	for _, test := range tests {
		result := IsEmpty(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestIsNotEmpty(t *testing.T) {
	assert.False(t, IsNotEmpty(""))
	assert.False(t, IsNotEmpty("   "))
	assert.True(t, IsNotEmpty("hello"))
}

func TestDefaultIfEmpty(t *testing.T) {
	assert.Equal(t, "default", DefaultIfEmpty("", "default"))
	assert.Equal(t, "default", DefaultIfEmpty("   ", "default"))
	assert.Equal(t, "hello", DefaultIfEmpty("hello", "default"))
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"hello", "olleh"},
		{"Hello, World!", "!dlroW ,olleH"},
		{"测试", "试测"},
	}

	for _, test := range tests {
		result := Reverse(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestTruncateWords(t *testing.T) {
	tests := []struct {
		input    string
		maxWords int
		expected string
	}{
		{"hello world test", 2, "hello world"},
		{"hello world test", 5, "hello world test"},
		{"hello world test", 0, ""},
		{"hello", 1, "hello"},
		{"", 2, ""},
	}

	for _, test := range tests {
		result := TruncateWords(test.input, test.maxWords)
		assert.Equal(t, test.expected, result)
	}
}

func TestTruncateChars(t *testing.T) {
	tests := []struct {
		input    string
		maxChars int
		suffix   string
		expected string
	}{
		{"hello world", 8, "...", "hello..."},
		{"hello world", 20, "...", "hello world"},
		{"hello world", 5, "...", "he..."},
		{"hello", 3, "...", "..."},
	}

	for _, test := range tests {
		result := TruncateChars(test.input, test.maxChars, test.suffix)
		assert.Equal(t, test.expected, result)
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello world", "helloWorld"},
		{"HelloWorld", "helloworld"},
		{"HELLO_WORLD", "helloWorld"},
		{"", ""},
	}

	for _, test := range tests {
		result := CamelCase(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"", ""},
	}

	for _, test := range tests {
		result := PascalCase(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"hello world", "hello_world"},
		{"hello-world", "hello_world"},
		{"helloWorld", "hello_world"},
		{"", ""},
	}

	for _, test := range tests {
		result := SnakeCase(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello-world"},
		{"hello world", "hello-world"},
		{"hello_world", "hello-world"},
		{"helloWorld", "hello-world"},
		{"", ""},
	}

	for _, test := range tests {
		result := KebabCase(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"HELLO", "Hello"},
		{"hELLO", "Hello"},
		{"", ""},
		{"h", "H"},
	}

	for _, test := range tests {
		result := Capitalize(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	assert.True(t, Contains(slice, "apple"))
	assert.True(t, Contains(slice, "banana"))
	assert.False(t, Contains(slice, "grape"))
	assert.False(t, Contains(slice, "Apple"))
}

func TestContainsIgnoreCase(t *testing.T) {
	slice := []string{"Apple", "Banana", "Cherry"}

	assert.True(t, ContainsIgnoreCase(slice, "apple"))
	assert.True(t, ContainsIgnoreCase(slice, "BANANA"))
	assert.False(t, ContainsIgnoreCase(slice, "grape"))
}

func TestRemoveDuplicates(t *testing.T) {
	input := []string{"a", "b", "c", "a", "b", "d"}
	expected := []string{"a", "b", "c", "d"}

	result := RemoveDuplicates(input)
	assert.Equal(t, expected, result)
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		input     string
		delimiter string
		expected  []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{" a , b , c ", ",", []string{"a", "b", "c"}},
		{"a,,b", ",", []string{"a", "b"}},
		{"", ",", []string{}},
		{"a", ",", []string{"a"}},
	}

	for _, test := range tests {
		result := SplitAndTrim(test.input, test.delimiter)
		assert.Equal(t, test.expected, result)
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"0", true},
		{"", false},
		{"12a", false},
		{"a12", false},
		{" 123", false},
	}

	for _, test := range tests {
		result := IsNumeric(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"ABC", true},
		{"", false},
		{"abc123", false},
		{"abc ", false},
	}

	for _, test := range tests {
		result := IsAlpha(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC", true},
		{"123", true},
		{"", false},
		{"abc ", false},
		{"abc-123", false},
	}

	for _, test := range tests {
		result := IsAlphaNumeric(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.org", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidEmail(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

func TestRandomString(t *testing.T) {
	result, err := RandomString(10)
	assert.NoError(t, err)
	assert.Len(t, result, 10)
	assert.True(t, IsAlphaNumeric(result))

	// Test that two calls produce different results
	result2, err := RandomString(10)
	assert.NoError(t, err)
	assert.NotEqual(t, result, result2)
}

func TestRandomAlphaString(t *testing.T) {
	result, err := RandomAlphaString(10)
	assert.NoError(t, err)
	assert.Len(t, result, 10)
	assert.True(t, IsAlpha(result))
}

func TestRandomNumericString(t *testing.T) {
	result, err := RandomNumericString(10)
	assert.NoError(t, err)
	assert.Len(t, result, 10)
	assert.True(t, IsNumeric(result))
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		padChar  rune
		expected string
	}{
		{"hello", 10, ' ', "     hello"},
		{"hello", 8, '0', "000hello"},
		{"hello", 3, ' ', "hello"},
		{"", 5, 'x', "xxxxx"},
	}

	for _, test := range tests {
		result := PadLeft(test.input, test.length, test.padChar)
		assert.Equal(t, test.expected, result)
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		padChar  rune
		expected string
	}{
		{"hello", 10, ' ', "hello     "},
		{"hello", 8, '0', "hello000"},
		{"hello", 3, ' ', "hello"},
		{"", 5, 'x', "xxxxx"},
	}

	for _, test := range tests {
		result := PadRight(test.input, test.length, test.padChar)
		assert.Equal(t, test.expected, result)
	}
}
