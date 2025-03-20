// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Attribute struct {
	Name  string
	Value string
}

// NewAttribute creates an Attribute from name and value byte slices, normalizing the key to lowercase and cleaning the
// value.
func newAttribute(name []byte, value []byte) Attribute {
	// Convert name to lowercase if needed.
	hasUpper := false
	for _, b := range name {
		if b >= 'A' && b <= 'Z' {
			hasUpper = true
			break
		}
	}

	var keyStr string
	if hasUpper {
		keyStr = string(bytes.ToLower(name))
	} else {
		keyStr = string(name)
	}

	// Fast path for simple values without newlines or comments.
	if !bytes.ContainsAny(value, "\n#") {
		return Attribute{
			Name:  keyStr,
			Value: string(bytes.TrimSpace(value)),
		}
	}

	// Process multi-line values or values with comments.
	var buf strings.Builder
	buf.Grow(len(value))
	startLine := 0
	inLine := false

	for i := 0; i < len(value); i++ {
		if value[i] == '\n' || i == len(value)-1 {
			// Handle the final character if it's not a newline.
			endPos := i
			if i == len(value)-1 && value[i] != '\n' {
				endPos = i + 1
			}

			line := value[startLine:endPos]

			// Handle line continuation.
			if inLine && len(line) > 0 && line[0] == '+' {
				line = line[1:]
			}

			// Handle comments.
			if commentIdx := bytes.IndexByte(line, '#'); commentIdx >= 0 {
				line = line[:commentIdx]
			}

			line = bytes.TrimSpace(line)
			if len(line) > 0 {
				if buf.Len() > 0 {
					buf.WriteByte(' ')
				}
				buf.Write(line)
			}

			startLine = i + 1
			inLine = true
		}
	}

	return Attribute{
		Name:  keyStr,
		Value: buf.String(),
	}
}

// parseAttributes parses the given buffer into a slice of Attributes.
func parseAttributes(buf []byte) ([]Attribute, error) {
	if len(buf) == 0 {
		return nil, errors.New("parseAttributes: object cannot be null")
	}

	// Pre-allocate attributes slice - typical RPSL objects have 5-15 attributes.
	attributes := make([]Attribute, 0, 16)
	pos := 0

	for pos < len(buf) {
		key, newPos, err := parseKey(buf, pos)
		if err != nil {
			return nil, err
		}

		pos = newPos
		value, newPos := parseValue(buf, pos)
		pos = newPos

		attributes = append(attributes, newAttribute(key, value))
	}

	return attributes, nil
}

// isValidKeyChar returns true if c is allowed in a key.
func isValidKeyChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '*'
}

// isLineContinuationChar returns true if c signifies a continuation of a line.
func isLineContinuationChar(c byte) bool {
	return c == ' ' || c == '\t' || c == '+'
}

// parseKey extracts a key ending at the first ':' and returns the key, the position after the colon, and an error if
// any.
func parseKey(buf []byte, pos int) ([]byte, int, error) {
	start := pos

	for pos < len(buf) {
		c := buf[pos]
		if c == ':' {
			if pos == start {
				return nil, 0, fmt.Errorf("parseKey: zero-sized key at pos %d", pos)
			}

			return buf[start:pos], pos + 1, nil
		}

		if !isValidKeyChar(c) {
			return nil, 0, fmt.Errorf("parseKey: illegal character '%c' at pos %d", c, pos)
		}

		pos++
	}

	return nil, 0, fmt.Errorf("parseKey: no key found starting at pos %d", start)
}

// parseValue extracts a value until a newline that is not followed by a continuation char.
func parseValue(buf []byte, pos int) ([]byte, int) {
	start := pos
	stop := pos

	for pos < len(buf) {
		c := buf[pos]
		pos++

		if c == '\r' {
			continue
		}

		if c == '\n' && pos < len(buf) {
			next := buf[pos]
			if isLineContinuationChar(next) {
				continue
			}

			break
		}

		stop = pos
	}

	return buf[start:stop], pos
}

// String returns a string representation of the Attribute.
func (a *Attribute) String() string {
	var str strings.Builder
	str.WriteString(a.Name)
	str.WriteString(":")
	str.WriteString(a.Value)
	return str.String()
}
