package oss

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type XmlDecoderLite struct {
	reader          io.Reader
	attributePrefix string
	useRawToken     bool
}

func NewXmlDecoderLite(r io.Reader) *XmlDecoderLite {
	return &XmlDecoderLite{
		reader:          r,
		attributePrefix: "+@",
		useRawToken:     true,
	}
}

func (dec *XmlDecoderLite) Decode(root *XmlNode) error {
	return dec.decodeXML(root)
}

type XmlNode struct {
	Children []*XmlChildren
	Data     []string
}

type XmlChildren struct {
	K string
	V []*XmlNode
}

func (n *XmlNode) addChild(s string, c *XmlNode) {
	if n.Children == nil {
		n.Children = make([]*XmlChildren, 0)
	}
	for _, childEntry := range n.Children {
		if childEntry.K == s {
			childEntry.V = append(childEntry.V, c)
			return
		}
	}
	n.Children = append(n.Children, &XmlChildren{K: s, V: []*XmlNode{c}})
}

func (n *XmlNode) value() any {
	if len(n.Children) > 0 {
		return n.GetMap()
	}
	if n.Data != nil {
		return n.Data[0]
	}
	return nil
}

func (n *XmlNode) GetMap() map[string]any {
	node := map[string]any{}
	for _, kv := range n.Children {
		label := kv.K
		children := kv.V
		if len(children) > 1 {
			vals := make([]any, 0)
			for _, child := range children {
				vals = append(vals, child.value())
			}
			node[label] = vals
		} else {
			node[label] = children[0].value()
		}
	}
	return node
}

type element struct {
	parent *element
	n      *XmlNode
	label  string
}

func (dec *XmlDecoderLite) decodeXML(root *XmlNode) error {
	xmlDec := xml.NewDecoder(dec.reader)

	started := false

	// Create first element from the root node
	elem := &element{
		parent: nil,
		n:      root,
	}

	getToken := func() (xml.Token, error) {
		if dec.useRawToken {
			return xmlDec.RawToken()
		}
		return xmlDec.Token()
	}

	for {
		t, e := getToken()
		if e != nil && !errors.Is(e, io.EOF) {
			return e
		}
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			elem = &element{
				parent: elem,
				n:      &XmlNode{},
				label:  se.Name.Local,
			}

			for _, a := range se.Attr {
				elem.n.addChild(dec.attributePrefix+a.Name.Local, &XmlNode{Data: []string{a.Value}})
			}
		case xml.CharData:
			newBit := trimNonGraphic(string(se))
			if !started && len(newBit) > 0 {
				return fmt.Errorf("invalid XML: Encountered chardata [%v] outside of XML node", newBit)
			}

			if len(newBit) > 0 {
				elem.n.Data = append(elem.n.Data, newBit)
			}
		case xml.EndElement:
			if elem.parent != nil {
				elem.parent.n.addChild(elem.label, elem.n)
			}
			elem = elem.parent
		}
		started = true
	}

	return nil
}

func trimNonGraphic(s string) string {
	if s == "" {
		return s
	}

	var first *int
	var last int
	for i, r := range []rune(s) {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			continue
		}

		if first == nil {
			f := i
			first = &f
			last = i
		} else {
			last = i
		}
	}

	if first == nil {
		return ""
	}

	return string([]rune(s)[*first : last+1])
}
