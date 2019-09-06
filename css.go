package css

import (
	"bytes"
)

const (
	space              = 32
	doubleQuote        = 34
	leftParenthesis    = 40
	rightParenthesis   = 41
	period             = 46
	colon              = 58
	leftSquareBracket  = 91
	rightSquareBracket = 93
	smallN             = 110
	smallO             = 111
	smallT             = 116
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
	// var d []string
	// for _, v := range r.Declarations {
	// 	d = append(d, v.String())
	// }
	return "" //strings.Join(d, ";")
}

// Selector define the elements to which a set of rules apply.
type Selector struct {
	Simple     Simple      `json:"simple"`
	Combinates []Combinate `json:"combinate,omitempty"`
}

func (v *Selector) encode(dst *bytes.Buffer) error {
	if err := v.Simple.encode(dst); err != nil {
		return err
	}

	if len(v.Combinates) > 0 {
		for _, i := range v.Combinates {
			if err := i.encode(dst); err != nil {
				return err
			}
		}
	}

	return nil
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
func (v *Simple) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Element); err != nil {
		return err
	}

	if len(v.Classes) > 0 {
		dst.WriteByte(period)
		dst.Write(bytes.Join(v.Classes, []byte{period}))
	}

	if len(v.Attributes) > 0 {
		for _, a := range v.Attributes {
			a.encode(dst)
		}
	}

	if len(v.PseudoElements) > 0 {
		for _, p := range v.PseudoElements {
			dst.Write([]byte{colon, colon})
			p.encode(dst)
		}
	}

	if len(v.PseudoClasses) > 0 {
		for _, p := range v.PseudoClasses {
			dst.WriteByte(colon)
			p.encode(dst)
		}
	}

	if len(v.Negations) > 0 {
		for _, s := range v.Negations {
			dst.Write([]byte{colon, smallN, smallO, smallT, leftParenthesis})
			s.encode(dst)
			dst.WriteByte(rightParenthesis)
		}
	}

	return nil
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

	if len(v.Func) > 0 {
		dst.WriteByte(40)
		if _, err := dst.Write(v.Func); err != nil {
			return err
		}
		dst.WriteByte(41)
	}

	return nil
}

// Element is a identity node or class name or #ID or universal
type Element struct {
	Value       string    `json:"value"`
	Attribute   Attribute `json:"attribute,omitempty"`
	PseudoClass string    `json:"pseudo_class,omitempty"`
}

func (v *Element) String() string {
	// var pc string

	// if len(v.PseudoClass) > 0 {
	// pc = "::" + v.PseudoClass
	// }

	return "" //v.Value + v.Attribute.String() + pc
}

// Attribute is a matcher of selector by attribute
type Attribute struct {
	Attr     []byte `json:"attr"`
	Operator []byte `json:"operator,omitempty"`
	Value    []byte `json:"value,omitempty"`
	Modifier []byte `json:"modifier,omitempty"`
}

func (v *Attribute) encode(dst *bytes.Buffer) error {
	dst.WriteByte(leftSquareBracket)

	if _, err := dst.Write(v.Attr); err != nil {
		return err
	}

	if _, err := dst.Write(v.Operator); err != nil {
		return err
	}

	if len(v.Value) > 0 {
		dst.WriteByte(doubleQuote)
		if _, err := dst.Write(v.Value); err != nil {
			return err
		}
		dst.WriteByte(doubleQuote)
	}

	if len(v.Modifier) > 0 {
		dst.WriteByte(space)
		if _, err := dst.Write(v.Modifier); err != nil {
			return err
		}
	}

	dst.WriteByte(rightSquareBracket)

	return nil
}

// Combinate the relationship between the selectors
type Combinate struct {
	Combinator []byte `json:"combinator"`
	Simple     Simple `json:"simple"`
}

func (v *Combinate) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Combinator); err != nil {
		return err
	}

	return v.Simple.encode(dst)
}

// Declaration is setting CSS properties
type Declaration struct {
	Property []byte   `json:"property"`
	Value    [][]byte `json:"value"`
}

func (v *Declaration) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Property); err != nil {
		return err
	}
	dst.WriteByte(colon)
	if _, err := dst.Write(bytes.Join(v.Value, []byte{space})); err != nil {
		return err
	}
	return nil
}
