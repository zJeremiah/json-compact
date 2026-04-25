package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type NodeKind int

const (
	NodeObject NodeKind = iota
	NodeArray
	NodeString
	NodeNumber
	NodeBool
	NodeNull
)

type Node struct {
	Kind      NodeKind
	Raw       string // raw literal for string/number/bool/null
	Items     []Node
	Entries   []Entry
	FlatWidth int // cached single-line width (-1 = not computed)
}

type Entry struct {
	Key   string
	Value Node
}

// flatWidth computes how many characters this node would take on a single line.
func flatWidth(n *Node) int {
	if n.FlatWidth >= 0 {
		return n.FlatWidth
	}
	w := computeFlatWidth(n)
	n.FlatWidth = w
	return w
}

func computeFlatWidth(n *Node) int {
	switch n.Kind {
	case NodeString, NodeNumber, NodeBool, NodeNull:
		return len(n.Raw)
	case NodeArray:
		if len(n.Items) == 0 {
			return 2 // []
		}
		// [ item, item, item ]
		w := 4 // "[ " + " ]"
		for i := range n.Items {
			w += flatWidth(&n.Items[i])
			if i < len(n.Items)-1 {
				w += 2 // ", "
			}
		}
		return w
	case NodeObject:
		if len(n.Entries) == 0 {
			return 2 // {}
		}
		// { "key": value, "key": value }
		w := 4 // "{ " + " }"
		for i := range n.Entries {
			e := &n.Entries[i]
			keyQuoted := `"` + e.Key + `"`
			w += len(keyQuoted) + 2 + flatWidth(&e.Value) // "key": value
			if i < len(n.Entries)-1 {
				w += 2 // ", "
			}
		}
		return w
	}
	return 0
}

// alignInfo holds per-key alignment widths for sibling objects.
type alignInfo struct {
	keyWidths   map[string]int // max width of each key (quoted)
	valueWidths map[string]int // max width of each value per key
}

// alignedInlineWidth computes how wide a single-line aligned rendering of
// this object would be, given the alignment info.
func alignedInlineWidth(n *Node, align *alignInfo) int {
	if n.Kind != NodeObject || len(n.Entries) == 0 {
		return flatWidth(n)
	}
	// "{ " + entries + " }"
	w := 4
	for i := range n.Entries {
		e := &n.Entries[i]
		kw := len(`"` + e.Key + `"`)
		if tw, ok := align.keyWidths[e.Key]; ok && tw > kw {
			kw = tw
		}
		vw := flatWidth(&e.Value)
		if tw, ok := align.valueWidths[e.Key]; ok && tw > vw {
			vw = tw
		}
		w += kw + 2 + vw // key: value
		if i < len(n.Entries)-1 {
			w += 2 // ", "
		}
	}
	return w
}

// computeAlignment checks if all children of a container are objects,
// and if so, returns alignment info for all keys that appear across them.
func computeAlignment(nodes []Node) *alignInfo {
	if len(nodes) < 2 {
		return nil
	}

	for i := range nodes {
		if nodes[i].Kind != NodeObject || len(nodes[i].Entries) == 0 {
			return nil
		}
	}

	ai := &alignInfo{
		keyWidths:   make(map[string]int),
		valueWidths: make(map[string]int),
	}
	for i := range nodes {
		for j := range nodes[i].Entries {
			e := &nodes[i].Entries[j]
			kw := len(`"` + e.Key + `"`)
			vw := flatWidth(&e.Value)
			if kw > ai.keyWidths[e.Key] {
				ai.keyWidths[e.Key] = kw
			}
			if vw > ai.valueWidths[e.Key] {
				ai.valueWidths[e.Key] = vw
			}
		}
	}
	return ai
}

// computeObjectValueAlignment checks if an object's values are all objects,
// used for map-of-objects alignment (like the tiles example).
func computeObjectValueAlignment(n *Node) *alignInfo {
	if n.Kind != NodeObject || len(n.Entries) < 2 {
		return nil
	}

	vals := make([]Node, len(n.Entries))
	for i, e := range n.Entries {
		vals[i] = e.Value
	}

	return computeAlignment(vals)
}

type formatter struct {
	buf        bytes.Buffer
	indent     int
	indentStr  string
	maxWidth   int
}

func newFormatter(indentSize, maxWidth int) *formatter {
	return &formatter{
		indentStr: strings.Repeat(" ", indentSize),
		maxWidth:  maxWidth,
	}
}

func (f *formatter) writeIndent() {
	for i := 0; i < f.indent; i++ {
		f.buf.WriteString(f.indentStr)
	}
}

func (f *formatter) currentIndentWidth() int {
	return f.indent * len(f.indentStr)
}

func (f *formatter) fitsOnLine(n *Node) bool {
	return f.currentIndentWidth()+flatWidth(n) <= f.maxWidth
}

func (f *formatter) formatNode(n *Node) {
	switch n.Kind {
	case NodeString, NodeNumber, NodeBool, NodeNull:
		f.buf.WriteString(n.Raw)

	case NodeArray:
		if len(n.Items) == 0 {
			f.buf.WriteString("[]")
			return
		}
		if f.fitsOnLine(n) {
			f.formatArrayInline(n, nil)
			return
		}

		align := computeAlignment(n.Items)

		f.buf.WriteString("[\n")
		f.indent++
		for i := range n.Items {
			f.writeIndent()
			if align != nil && n.Items[i].Kind == NodeObject &&
				f.currentIndentWidth()+alignedInlineWidth(&n.Items[i], align) <= f.maxWidth {
				f.formatObjectInlineAligned(&n.Items[i], align)
			} else {
				f.formatNode(&n.Items[i])
			}
			if i < len(n.Items)-1 {
				f.buf.WriteString(", ")
			}
			f.buf.WriteString("\n")
		}
		f.indent--
		f.writeIndent()
		f.buf.WriteString("]")

	case NodeObject:
		if len(n.Entries) == 0 {
			f.buf.WriteString("{}")
			return
		}
		if f.fitsOnLine(n) {
			f.formatObjectInline(n, nil)
			return
		}

		valAlign := computeObjectValueAlignment(n)

		maxKeyW := 0
		for _, e := range n.Entries {
			kw := len(`"` + e.Key + `"`)
			if kw > maxKeyW {
				maxKeyW = kw
			}
		}

		f.buf.WriteString("{\n")
		f.indent++
		for i := range n.Entries {
			e := &n.Entries[i]
			f.writeIndent()

			keyQuoted := `"` + e.Key + `"`
			for len(keyQuoted) < maxKeyW {
				keyQuoted += " "
			}

			if valAlign != nil && e.Value.Kind == NodeObject {
				valInlineW := alignedInlineWidth(&e.Value, valAlign)
				lineW := f.currentIndentWidth() + maxKeyW + 2 + valInlineW

				if lineW <= f.maxWidth {
					f.buf.WriteString(keyQuoted)
					f.buf.WriteString(": ")
					f.formatObjectInlineAligned(&e.Value, valAlign)
				} else {
					f.buf.WriteString(keyQuoted)
					f.buf.WriteString(": ")
					f.formatNode(&e.Value)
				}
			} else {
				f.buf.WriteString(keyQuoted)
				f.buf.WriteString(": ")
				f.formatNode(&e.Value)
			}

			if i < len(n.Entries)-1 {
				f.buf.WriteString(", ")
			}
			f.buf.WriteString("\n")
		}
		f.indent--
		f.writeIndent()
		f.buf.WriteString("}")
	}
}

func (f *formatter) formatArrayInline(n *Node, align *alignInfo) {
	f.buf.WriteString("[ ")
	for i := range n.Items {
		if align != nil && n.Items[i].Kind == NodeObject {
			f.formatObjectInlineAligned(&n.Items[i], align)
		} else {
			f.formatNode(&n.Items[i])
		}
		if i < len(n.Items)-1 {
			f.buf.WriteString(", ")
		}
	}
	f.buf.WriteString(" ]")
}

func (f *formatter) formatObjectInline(n *Node, align *alignInfo) {
	f.buf.WriteString("{ ")
	for i := range n.Entries {
		e := &n.Entries[i]
		f.buf.WriteString(`"` + e.Key + `"`)
		f.buf.WriteString(": ")
		f.formatNode(&e.Value)
		if i < len(n.Entries)-1 {
			f.buf.WriteString(", ")
		}
	}
	f.buf.WriteString(" }")
}

// formatObjectInlineAligned renders an object on a single line with aligned fields.
func (f *formatter) formatObjectInlineAligned(n *Node, align *alignInfo) {
	f.buf.WriteString("{ ")
	for i := range n.Entries {
		e := &n.Entries[i]
		keyQuoted := `"` + e.Key + `"`

		// Pad key if needed
		if targetW, ok := align.keyWidths[e.Key]; ok && targetW > 0 {
			for len(keyQuoted) < targetW {
				keyQuoted += " "
			}
		}

		f.buf.WriteString(keyQuoted)
		f.buf.WriteString(": ")

		// Pad value
		valStr := nodeToString(&e.Value)
		if targetW, ok := align.valueWidths[e.Key]; ok {
			valStr = padValue(&e.Value, valStr, targetW)
		}
		f.buf.WriteString(valStr)

		if i < len(n.Entries)-1 {
			f.buf.WriteString(", ")
		}
	}
	f.buf.WriteString(" }")
}

func nodeToString(n *Node) string {
	switch n.Kind {
	case NodeString, NodeNumber, NodeBool, NodeNull:
		return n.Raw
	default:
		tmp := newFormatter(1, 9999)
		tmp.formatNode(n)
		return tmp.buf.String()
	}
}

// padValue pads a value string to targetWidth.
// Numbers are right-aligned (left-padded), strings are left-aligned (right-padded).
func padValue(n *Node, valStr string, targetWidth int) string {
	if len(valStr) >= targetWidth {
		return valStr
	}
	pad := targetWidth - len(valStr)
	switch n.Kind {
	case NodeNumber:
		return strings.Repeat(" ", pad) + valStr
	default:
		return valStr + strings.Repeat(" ", pad)
	}
}

func run() error {
	filePath := flag.String("f", "", "read JSON from file instead of stdin")
	maxWidth := flag.Int("w", 120, "max line width before expanding objects")
	indentSize := flag.Int("indent", 1, "number of spaces per indent level")
	flag.Parse()

	var input []byte
	var err error

	if *filePath != "" {
		input, err = os.ReadFile(*filePath)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
	}

	if len(bytes.TrimSpace(input)) == 0 {
		return fmt.Errorf("empty input")
	}

	if !json.Valid(input) {
		return fmt.Errorf("invalid JSON")
	}

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()
	root, err := parseValueNumber(dec)
	if err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	fmtr := newFormatter(*indentSize, *maxWidth)
	fmtr.formatNode(&root)
	fmtr.buf.WriteString("\n")

	os.Stdout.Write(fmtr.buf.Bytes())
	return nil
}

// parseValueNumber is like parseValue but uses json.Number to preserve
// integer formatting.
func parseValueNumber(dec *json.Decoder) (Node, error) {
	tok, err := dec.Token()
	if err != nil {
		return Node{}, err
	}

	switch v := tok.(type) {
	case json.Delim:
		switch v {
		case '{':
			n := Node{Kind: NodeObject, FlatWidth: -1}
			for dec.More() {
				keyTok, err := dec.Token()
				if err != nil {
					return Node{}, err
				}
				key := keyTok.(string)
				val, err := parseValueNumber(dec)
				if err != nil {
					return Node{}, err
				}
				n.Entries = append(n.Entries, Entry{Key: key, Value: val})
			}
			if _, err := dec.Token(); err != nil {
				return Node{}, err
			}
			return n, nil
		case '[':
			n := Node{Kind: NodeArray, FlatWidth: -1}
			for dec.More() {
				val, err := parseValueNumber(dec)
				if err != nil {
					return Node{}, err
				}
				n.Items = append(n.Items, val)
			}
			if _, err := dec.Token(); err != nil {
				return Node{}, err
			}
			return n, nil
		}
	case string:
		raw, _ := json.Marshal(v)
		return Node{Kind: NodeString, Raw: string(raw), FlatWidth: -1}, nil
	case json.Number:
		return Node{Kind: NodeNumber, Raw: v.String(), FlatWidth: -1}, nil
	case bool:
		if v {
			return Node{Kind: NodeBool, Raw: "true", FlatWidth: -1}, nil
		}
		return Node{Kind: NodeBool, Raw: "false", FlatWidth: -1}, nil
	case nil:
		return Node{Kind: NodeNull, Raw: "null", FlatWidth: -1}, nil
	}

	return Node{}, fmt.Errorf("unexpected token: %v", tok)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "json-compact: %s\n", err)
		os.Exit(1)
	}
}
