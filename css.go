package css

import (
	"bytes"
	"strings"
)

// Statements sets Statement
type Statements []Statement

// Statement is a building block
type Statement struct {
	Identifier  string            `json:"type"`
	Information interface{}       `json:"information,omitempty"`
	Nested      *StatementRuleset `json:"nested,omitempty"`
}

// StatementRuleset is nested block
type StatementRuleset struct {
	Statement Statement `json:"statement"`
	Ruleset   Ruleset   `json:"ruleset"`
}

// Ruleset is a collection of CSS declarations
type Ruleset struct {
	Selectors    []Selector    `json:"selectors"`
	Declarations []Declaration `json:"declarations"`
}

func (r *Ruleset) String() string {
	var d []string
	for _, v := range r.Declarations {
		d = append(d, v.String())
	}
	return strings.Join(d, ";")
}

// Selector define the elements to which a set of rules apply.
type Selector struct {
	Element Element `json:"element,omitempty"`

	PseudoElement string `json:"pseudo_element,omitempty"`

	Simple Simple `json:"simple"`
}

// Simple is a simple selector
type Simple struct {
	Element        []byte      `json:"element,omitempty"`
	Classes        [][]byte    `json:"classes,omitempty"`
	Attributes     []Attribute `json:"attributes,omitempty"`
	PseudoElements []Pseudo    `json:"pseudo_elements,omitempty"`
	PseudoClasses  []Pseudo    `json:"pseudo_classes,omitempty"`
	Negations      []Simple    `json:"negations,omitempty"`
}

// Encode to CSS
func (v *Simple) Encode() []byte {
	var ret []byte

	ret = append(ret, v.Element...)
	if len(v.Classes) > 0 {
		ret = append(ret, []byte(".")...)
		ret = append(ret, bytes.Join(v.Classes, []byte("."))...)
	}

	if len(v.Attributes) > 0 {
		// for _, a := range v.Attributes {
		// ret += a.String()
		// }
	}

	if len(v.PseudoElements) > 0 {
		// for _, p := range v.PseudoElements {
		// 	ret = append(ret, []byte("::")...)
		// 	// ret = append(ret, p.Encode()...)
		// }
	}

	if len(v.PseudoClasses) > 0 {
		// for _, p := range v.PseudoClasses {
		// 	ret = append(ret, []byte(":")...)
		// 	// ret = append(ret, p.Encode()...)
		// }
	}

	if len(v.Negations) > 0 {
		for _, s := range v.Negations {
			ret = append(ret, []byte(":not(")...)
			ret = append(ret, s.Encode()...)
			ret = append(ret, []byte(")")...)
		}
	}

	return ret
}

// Pseudo is a pseudo-class
type Pseudo struct {
	Ident []byte `json:"ident,omitempty"`
	Func  []byte `json:"func,omitempty"`
}

// Encode to CSS
func (v *Pseudo) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Ident); err != nil {
		return err
	}

	dst.WriteByte(40)

	if _, err := dst.Write(v.Func); err != nil {
		return err
	}

	dst.WriteByte(41)

	return nil
}

// Element is a identity node or class name or #ID or universal
type Element struct {
	Value       string    `json:"value"`
	Attribute   Attribute `json:"attribute,omitempty"`
	PseudoClass string    `json:"pseudo_class,omitempty"`
}

func (v *Element) String() string {
	var pc string

	if len(v.PseudoClass) > 0 {
		pc = "::" + v.PseudoClass
	}

	return v.Value + v.Attribute.String() + pc
}

// Combinator the relationship between the selectors
type Combinator struct {
	First    Selector `json:"first"`
	Operator string   `json:"operator"`
	Second   Selector `json:"second"`
}

// Attribute is a matcher of selector by attribute
type Attribute struct {
	Attr     string `json:"attr"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
	Modifier string `json:"modifier,omitempty"`
}

func (v *Attribute) String() string {
	var m string
	if len(v.Modifier) > 0 {
		m = " " + v.Modifier
	}
	return "[" + v.Attr + v.Operator + v.Value + m + "]"
}

// Declaration is setting CSS properties
type Declaration struct {
	Property string   `json:"property"`
	Value    []string `json:"value"`
}

func (d *Declaration) String() string {
	return d.Property + ":" + strings.Join(d.Value, " ")
}
