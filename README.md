# obsidian-pipeline

# Overview
`obsidian-pipeline` or `obp` enables headless management of your Obsidian vault.

`obp` is primarily targeted at users who are programmatically interacting
with their vaults, publishing them via static site generators or doing
quality control with scheduled tasks.

# Features
## Validation
Markdown allows you to preface your document with another, typically in
YAML, that contains metadata about the document. A common use-case for this is
in static site generators; Hugo checks your Markdown document's frontmatter for
properties like `title` to help it make rendering decisions.

When you have data, making sure it's consistent is pretty important. For example,
Hugo uses the boolean `draft` property to determine whether a post will be
included in the generated site output. A personal blog is a low
stakes example, but a business managing internal vs external documentation
or otherwise enforcing standards could save themselves a lot of effort and
headache if it were possible to automate the process of verifying the layout
and disposition of their data.

The good news is that we don't have to invent anything here. Schemas exist
precisely for this reason
