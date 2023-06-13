package obp_test

import (
	"bytes"
	"testing"

	"code.ndumas.com/ndumas/obsidian-pipeline"
)

func Test_BasicValidation(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		b        *bytes.Buffer
		expected error
	}{
		{
			name: "KeyMissing",
			b: bytes.NewBufferString(`
---
boop: "bop"
---
# Markdown Content
`),
			expected: nil,
		},
		{
			name: "KeyTypeMismatch",
			b: bytes.NewBufferString(`
---
title: 2
---
# Markdown Content
`),
			expected: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := obp.Validate("https://schemas.ndumas.com/obsidian/note.schema.json", tc.b)
			if err == nil {
				t.Log("Expected Validate() to fail on input")
				t.Fail()
			}
		})
	}
}
