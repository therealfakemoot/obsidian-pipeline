# obsidian-pipeline

# Overview
`obsidian-pipeline` or `obp` enables headless management of your Obsidian vault.

`obp` is primarily targeted at users who are programmatically interacting
with their vaults, publishing them via static site generators or doing
quality control with scheduled tasks.

# Features
`obp` is in pre-release status for now, so the featureset is limited.

```
obp-linux help
a suite of tools for managing your obsidian vault

Usage:
  obp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  validate    loads a note and ensures its frontmatter follows the provided protobuf schema

Flags:
      --config string   config file (default "~/.obp.toml")
  -h, --help            help for obp

Additional help topics:
  obp hugo       convert a set of Obsidian notes into a Hugo compatible directory structure

Use "obp [command] --help" for more information about a command.

```
## Validation

Using the validation feature requires that the JSONSchema files are hosted on
some HTTP server somewhere. I did this to simplify the code and avoid resolving
relative pathing along the dependency tree.

You can find some pre-made schemas [here](https://schemas.ndumas.com), a site
powered by my [json-schemas](https://github.com/therealfakemoot/json-schemas) repository. If you
would like to submit a schema that you've designed.

### Usage
```
Usage:
  obp validate [flags]

Flags:
      --format string   output format [markdown, json, csv] (default "markdown")
  -h, --help            help for validate
  -s, --schema string   path to protobuf file (default "base.schema")
  -t, --target string   directory containing validation targets

Global Flags:
      --config string   config file (default "~/.obp.toml")
```

Markdown output can be piped into [glow](https://github.com/charmbracelet/glow) for a little
extra pizazz, but JSON is available for programmatic handling.
```
./bin/obp-linux validate --format markdown -s https://schemas.ndumas.com/obsidian/note.schema.json -t Resou
rces/blog/published/|glow

   Validation Errors for "/home/ndumas/work/obsidian-pipeline/Resources/blog/published/schema-bad-tags.md"

          VALIDATION RULE       │ FAILING PROPERTY │             ERROR
  ──────────────────────────────┼──────────────────┼─────────────────────────────────
    /properties/tags/items/type │ /tags/2          │ expected string, but got
                                │                  │ number

   Validation Errors for "/home/ndumas/work/obsidian-pipeline/Resources/blog/published/schema-bad.md"

          VALIDATION RULE        │ FAILING PROPERTY │             ERROR
  ───────────────────────────────┼──────────────────┼─────────────────────────────────
    /properties/description/type │ /description     │ expected string, but got
                                 │                  │ number
    /properties/tags/type        │ /tags            │ expected array, but got string
```
