# json-compact

A CLI tool that reformats JSON into a compact, human-readable form. Small objects are collapsed onto single lines, and sibling values are column-aligned for easy scanning.

## Install

```bash
go install github.com/zjeremiah/json-compact@latest
```

Or build from source:

```bash
go build -o json-compact .
```

## Usage

```
json-compact [flags]
```

Reads JSON from **stdin** by default and writes the compacted output to **stdout**.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-f <file>` | _(stdin)_ | Read JSON from a file instead of stdin |
| `-w <int>` | `120` | Max line width before an object is expanded to multiple lines |
| `-indent <int>` | `1` | Number of spaces per indent level |

### Examples

```bash
# Pipe from stdin
cat data.json | json-compact

# Read from a file
json-compact -f data.json

# Narrow output width
json-compact -w 80 -f data.json

# Use 2-space indentation
json-compact -indent 2 -f data.json
```

### Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Invalid JSON, parse error, or file read error |

## What it does

Given this input:

```json
{
  "animations": {
    "40": [
      {
        "duration": 150,
        "tileid": 40
      },
      {
        "duration": 150,
        "tileid": 521
      },
      {
        "duration": 150,
        "tileid": 41
      },
      {
        "duration": 150,
        "tileid": 522
      }
    ]
  }
}
```

`json-compact` produces:

```json
{
 "animations": {
  "40": [
   { "duration": 150, "tileid":  40 },
   { "duration": 150, "tileid": 521 },
   { "duration": 150, "tileid":  41 },
   { "duration": 150, "tileid": 522 }
  ]
 }
}
```

### Formatting rules

- **Single-line collapse** -- Objects and arrays that fit within the max line width (including indentation) are rendered on one line.
- **Column alignment** -- When sibling objects share the same set of keys (e.g. elements of an array, or values of a map), their keys and values are padded to align into columns:
  - String values are left-aligned (right-padded)
  - Number values are right-aligned (left-padded)
  - Object keys of differing lengths are right-padded so colons align
- **Key padding** -- In any expanded object, keys are padded to the same width so all colons line up.
- **Key order preserved** -- Keys appear in the same order as the input.
- **Width-aware** -- Aligned inline objects that would exceed the max line width are expanded to multiple lines instead.

### Map-of-objects example

Input:

```json
{"height":16,"tiles":{"0,0":{"id":1578,"tileset":"Objects"},"0,10":{"id":1506,"tileset":"Objects"},"4,1":{"id":837,"tileset":"Objects_animated"},"4,2":{"id":838,"tileset":"Objects_animated"}},"width":16,"x":16,"y":0}
```

Output:

```json
{
 "height": 16,
 "tiles" : {
  "0,0" : { "id": 1578, "tileset": "Objects"          },
  "0,10": { "id": 1506, "tileset": "Objects"          },
  "4,1" : { "id":  837, "tileset": "Objects_animated" },
  "4,2" : { "id":  838, "tileset": "Objects_animated" }
 },
 "width" : 16,
 "x"     : 16,
 "y"     : 0
}
```

## No dependencies

Built entirely with the Go standard library. No external modules required.
