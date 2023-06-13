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
		{
			name: "GoodSchema",
			b: bytes.NewBufferString(`
---
draft: false
title: "Mapping Aardwolf with Graphviz and Golang"
aliases: ["Mapping Aardwolf with Graphviz"]
series: ["mapping-aardwolf"]
date: "2023-04-06"
author: "Nick Dumas"
cover: ""
keywords: [""]
description: "Maxing out your CPU for fun and profit with dense graphs, or how I'm attempting to follow through on my plan to work on projects with more visual
 outputs"
showFullContent: false
tags:
- graphviz
- graph
- aardwolf
- golang
---

## Textual Cartography
Aardwolf has a fairly active developer community, people who write and maintain plugins and try to map the game world and its contents.
`),
			expected: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := obp.Validate("https://schemas.ndumas.com/obsidian/note.schema.json", tc.b)
			if err == tc.expected {
				t.Log("Expected Validate() to fail on input")
				t.Fail()
			}
		})
	}
}
