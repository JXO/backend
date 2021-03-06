// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package loaders

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/jxo/lime/loaders/plist"
	"github.com/jxo/lime/parser"
)

func plistconv(buf *bytes.Buffer, node *parser.Node) error {
	switch node.Name {
	case "Key":
		buf.WriteString("\"" + node.Data() + "\": ")
	case "String":
		n := node.Data()
		n = strings.Replace(n, "\\", "\\\\", -1)
		n = strings.Replace(n, "\"", "\\\"", -1)
		n = strings.Replace(n, "\n", "\\n", -1)
		n = strings.Replace(n, "\t", "\\t", -1)
		n = strings.Replace(n, "&gt;", ">", -1)
		n = strings.Replace(n, "&lt;", "<", -1)
		buf.WriteString("\"" + n + "\"")
	case "EmptyString":
		buf.WriteString("\"\"")
	case "Dictionary":
		buf.WriteString("{\n\t")
		for i, child := range node.Children {
			if i != 0 && i&1 == 0 {
				buf.WriteString(",\n\t")
			}
			if err := plistconv(buf, child); err != nil {
				return err
			}
		}

		buf.WriteString("}\n")
	case "Array":
		buf.WriteString("[\n\t")
		for i, child := range node.Children {
			if i != 0 {
				buf.WriteString(",\n\t")
			}

			if err := plistconv(buf, child); err != nil {
				return err
			}
		}

		buf.WriteString("]\n\t")
	case "Integer", "Bool":
		buf.WriteString(node.Data())
	case "EndOfFile":
	default:
		return errors.New(fmt.Sprintf("Unhandled node: %s", node.Name))
	}
	return nil
}

func LoadPlist(data []byte, intf interface{}) error {
	var p plist.PLIST
	r := strings.NewReplacer("\r", "", "\v", "")
	if !p.Parse(r.Replace(string(data))) {
		return p.Error()
	}
	var (
		root = p.RootNode()
		buf  bytes.Buffer
	)
	for _, child := range root.Children {
		if err := plistconv(&buf, child); err != nil {
			return err
		}
	}
	return LoadJSON(buf.Bytes(), intf)
}
