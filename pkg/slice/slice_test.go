package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendIfDNE(t *testing.T) {
	testCases := []struct {
		Input          []string
		StringToAppend string
		ExpectedOutput []string
	}{
		{[]string{"foo"}, "bar", []string{"foo", "bar"}},
		{[]string{"foo", "bar"}, "foo", []string{"foo", "bar"}},
		{[]string{"foo", "bar"}, "bar", []string{"foo", "bar"}},
		{[]string{"foo", "bar"}, "test", []string{"foo", "bar", "test"}},
	}

	for _, testCase := range testCases {
		input := testCase.Input
		stringToAppend := testCase.StringToAppend
		output := AppendIfDNE(input, stringToAppend)

		assert.Equal(t, testCase.Input, input)
		assert.Equal(t, testCase.ExpectedOutput, output)
	}
}

func TestRemoveFromStringSlice(t *testing.T) {
	testCases := []struct {
		Input          []string
		StringToRemove string
		ExpectedOutput []string
	}{
		{[]string{"foo"}, "foo", []string{}},
		{[]string{"foo", "foo"}, "foo", []string{}},
		{[]string{"foo", "foo", "foo"}, "foo", []string{}},
		{[]string{"foo", "bar", "foo"}, "foo", []string{"bar"}},
		{[]string{"foo", "bar", "foo", "bar", "foo"}, "foo", []string{"bar", "bar"}},
		{[]string{"bar"}, "foo", []string{"bar"}},
	}

	for _, testCase := range testCases {
		input := testCase.Input
		stringToRemove := testCase.StringToRemove
		output := RemoveFromStringSlice(input, stringToRemove)

		assert.Equal(t, testCase.Input, input)
		assert.Equal(t, testCase.ExpectedOutput, output)
	}
}

func TestUniqueStringSlice(t *testing.T) {
	testCases := []struct {
		Input          []string
		ExpectedOutput []string
	}{
		{[]string{"foo", "bar"}, []string{"foo", "bar"}},
		{[]string{"foo", "bar", "bar"}, []string{"foo", "bar"}},
		{[]string{"foo", "foo", "bar", "bar"}, []string{"foo", "bar"}},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.ExpectedOutput, UniqueStringSlice(testCase.Input))
	}
}
