package css2json

import (
	"bytes"
	"encoding/json"
	"errors"
)

const (
	space              = 32
	doubleQuote        = 34
	leftParenthesis    = 40
	rightParenthesis   = 41
	comma              = 44
	period             = 46
	colon              = 58
	semicolon          = 59
	atSign             = 64
	leftSquareBracket  = 91
	rightSquareBracket = 93
	smallN             = 110
	smallO             = 111
	smallT             = 116
	leftCurlyBracket   = 123
	rightCurlyBracket  = 125
)

var (
	// ErrNotExistsDeclaration
	ErrNotExistsDeclaration = errors.New("not exists declaration")
	// ErrNotExistsTypeIdentifier
	ErrNotExistsTypeIdentifier = errors.New("not exists type of identifier")
)

// Statements sets Statement
type Statements []Statement

// Encode statements to CSS
func Encode(s Statements) ([]byte, error) {
	buf := &bytes.Buffer{}

	for _, i := range s {
		if err := i.encode(buf); err != nil {
			return nil, err
		}
		buf.WriteByte(semicolon)
	}

	return buf.Bytes(), nil
}

// Statement is a building block
type Statement struct {
	AtRule  *AtRule  `json:"atrule,omitempty"`
	Ruleset *Ruleset `json:"ruleset,omitempty"`
}

func (v *Statement) encode(dst *bytes.Buffer) error {
	if v.AtRule != nil {
		if err := v.AtRule.encode(dst); err != nil {
			return err
		}
	}
	if v.Ruleset != nil {
		if err := v.Ruleset.encode(dst); err != nil {
			return err
		}
	}
	return nil
}

// AtRule https://developer.mozilla.org/en-US/docs/Web/CSS/At-rule
type AtRule struct {
	Identifier Identifier   `json:"ident"`
	Nested     []*Statement `json:"nested,omitempty"`
}

func (v *AtRule) encode(dst *bytes.Buffer) error {
	if err := v.Identifier.encode(dst); err != nil {
		return err
	}

	if v.Nested != nil {
		dst.WriteByte(leftCurlyBracket)
		for _, i := range v.Nested {
			if err := i.encode(dst); err != nil {
				return err
			}
		}
		dst.WriteByte(rightCurlyBracket)
	}

	return nil
}

// Ruleset is a collection of CSS declarations
type Ruleset struct {
	Selectors    []Selector    `json:"selectors"`
	Declarations []Declaration `json:"declarations"`
}

func (v *Ruleset) encode(dst *bytes.Buffer) error {
	if len(v.Declarations) == 0 {
		return ErrNotExistsDeclaration
	}

	for idx, s := range v.Selectors {
		if err := s.encode(dst); err != nil {
			return err
		}
		if len(v.Selectors)-1 > idx {
			dst.WriteByte(comma)
		}
	}

	dst.WriteByte(leftCurlyBracket)
	for idx, d := range v.Declarations {
		if err := d.encode(dst); err != nil {
			return err
		}
		if len(v.Declarations)-1 > idx {
			dst.WriteByte(semicolon)
		}
	}
	dst.WriteByte(rightCurlyBracket)

	return nil
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
	Element        TextBytes   `json:"element,omitempty"`
	Classes        []TextBytes `json:"classes,omitempty"`
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
		dst.Write(bytes.Join(sliceValuesRaw(v.Classes), []byte{period}))
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
	Ident TextBytes `json:"ident,omitempty"`
	Func  TextBytes `json:"func,omitempty"`
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

// Attribute is a matcher of selector by attribute
type Attribute struct {
	Attr     TextBytes `json:"attr"`
	Operator TextBytes `json:"operator,omitempty"`
	Value    TextBytes `json:"value,omitempty"`
	Modifier TextBytes `json:"modifier,omitempty"`
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
	Combinator TextBytes `json:"combinator"`
	Simple     Simple    `json:"simple"`
}

func (v *Combinate) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Combinator); err != nil {
		return err
	}

	return v.Simple.encode(dst)
}

// Value of property
type Value struct {
	ValueSpace []TextBytes `json:"values,omitempty"`
}

func (v *Value) encode(dst *bytes.Buffer) error {
	data := bytes.Join(sliceValuesRaw(v.ValueSpace), []byte{space})
	if _, err := dst.Write(data); err != nil {
		return err
	}
	return nil
}

func sliceValuesRaw(v []TextBytes) [][]byte {
	ret := make([][]byte, len(v))
	for k, i := range v {
		ret[k] = i
	}
	return ret
}

// Declaration is setting CSS properties
type Declaration struct {
	Property TextBytes `json:"property"`
	Values   []Value   `json:"values,omitempty"`
}

func (v *Declaration) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Property); err != nil {
		return err
	}
	dst.WriteByte(colon)

	for idx, i := range v.Values {
		if err := i.encode(dst); err != nil {
			return err
		}
		if len(v.Values)-1 > idx {
			dst.WriteByte(comma)
		}
	}

	return nil
}

// TextBytes is a hack to get JSON to emit a []byte as a string
type TextBytes []byte

// MarshalJSON marshal TextBytes
func (v TextBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(v))
}

// UnmarshalJSON unmarshal TextBytes
func (v *TextBytes) UnmarshalJSON(b []byte) error {
	var a string
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*v = TextBytes(a)

	return nil
}

type encoder interface {
	encode(*bytes.Buffer) error
}

type writer func(*bytes.Buffer) error

func encodeItemsIfExists(items []encoder, dst *bytes.Buffer, before, after writer) error {
	if len(items) <= 0 {
		return nil
	}
	for _, i := range items {
		if before != nil {
			if err := before(dst); err != nil {
				return err
			}
		}
		if err := i.encode(dst); err != nil {
			return err
		}
		if after != nil {
			if err := after(dst); err != nil {
				return err
			}
		}
	}
	return nil
}

// Identifier is a at-rule
type Identifier struct {
	Type        TextBytes   `json:"type"`
	Information Information `json:"info,omitempty"`
}

// UnmarshalJSON unmarshal TextBytes
func (v *Identifier) UnmarshalJSON(b []byte) error {
	var (
		info  Information
		stuff map[string]interface{}
	)

	if err := json.Unmarshal(b, &stuff); err != nil {
		return err
	}

	tmp, err := json.Marshal(stuff["info"])
	if err != nil {
		return err
	}

	ident := stuff["type"].(string)
	if _, ok := identifierTypes[ident]; !ok {
		return ErrNotExistsTypeIdentifier
	}
	info = identifierTypes[ident].(encoder)

	if err := json.Unmarshal(tmp, info); err != nil {
		return err
	}
	v.Information = info

	return nil
}

func (v *Identifier) encode(dst *bytes.Buffer) error {
	if _, ok := identifierTypes[string(v.Type)]; !ok {
		return ErrNotExistsTypeIdentifier
	}

	dst.WriteByte(atSign)

	if _, err := dst.Write(v.Type); err != nil {
		return err
	}

	if v.Information != nil {
		dst.WriteByte(space)
	}

	if err := v.Information.encode(dst); err != nil {
		return err
	}

	return nil
}

// Information is a information of at-rule
type Information interface {
	encoder
}

var identifierTypes = map[string]interface{}{
	"charset":   &CharsetInformation{},
	"keyframes": &KeyframesInformation{},
	"media":     &MediaInformation{},
	"font-face": &FontFaceInformation{},
}

// CharsetInformation https://developer.mozilla.org/en-US/docs/Web/CSS/@charset
type CharsetInformation struct {
	Value TextBytes `json:"value"`
}

func (v *CharsetInformation) encode(dst *bytes.Buffer) error {
	dst.WriteByte(doubleQuote)
	if _, err := dst.Write(v.Value); err != nil {
		return err
	}
	dst.WriteByte(doubleQuote)

	return nil
}

// KeyframesInformation https://developer.mozilla.org/en-US/docs/Web/CSS/@keyframes
type KeyframesInformation struct {
	Value TextBytes `json:"value"`
}

func (v *KeyframesInformation) encode(dst *bytes.Buffer) error {
	if _, err := dst.Write(v.Value); err != nil {
		return err
	}

	return nil
}

// MediaInformation https://developer.mozilla.org/en-US/docs/Web/CSS/@media
type MediaInformation struct {
	Queries []Query `json:"queries"`
}

func (v *MediaInformation) encode(dst *bytes.Buffer) error {
	if len(v.Queries) == 0 {
		return nil
	}

	for idx, i := range v.Queries {
		if err := i.encode(dst); err != nil {
			return err
		}
		if len(v.Queries)-1 > idx {
			dst.WriteByte(comma)
		}
	}

	return nil
}

// Query is a media-condition
type Query struct {
	Type       *Type       `json:"type,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

func (v *Query) encode(dst *bytes.Buffer) error {
	if v.Type != nil {
		if err := v.Type.encode(dst); err != nil {
			return err
		}
	}

	if len(v.Conditions) > 0 {
		for _, i := range v.Conditions {
			if v.Type != nil {
				dst.WriteByte(space)
			}
			if err := i.encode(dst); err != nil {
				return err
			}
		}
	}

	return nil
}

// Type is a media type
type Type struct {
	Operator TextBytes `json:"operator,omitempty"`
	Value    TextBytes `json:"value"`
}

func (v *Type) encode(dst *bytes.Buffer) error {
	if len(v.Operator) > 0 {
		if _, err := dst.Write(v.Operator); err != nil {
			return err
		}
		dst.WriteByte(space)
	}

	if _, err := dst.Write(v.Value); err != nil {
		return err
	}

	return nil
}

// Condition is a media condition
type Condition struct {
	Operator TextBytes `json:"operator,omitempty"`
	Feature  TextBytes `json:"feature"`
	Value    TextBytes `json:"value"`
}

func (v *Condition) encode(dst *bytes.Buffer) error {
	if len(v.Operator) > 0 {
		if _, err := dst.Write(v.Operator); err != nil {
			return err
		}
		dst.WriteByte(space)
	}

	dst.WriteByte(leftParenthesis)

	if _, err := dst.Write(v.Feature); err != nil {
		return err
	}

	dst.WriteByte(colon)

	if _, err := dst.Write(v.Value); err != nil {
		return err
	}

	dst.WriteByte(rightParenthesis)

	return nil
}

// FontFaceInformation https://developer.mozilla.org/en-US/docs/Web/CSS/@font-face
type FontFaceInformation struct {
	Declarations []Declaration `json:"declarations"`
}

func (v *FontFaceInformation) encode(dst *bytes.Buffer) error {
	dst.WriteByte(leftCurlyBracket)
	for idx, i := range v.Declarations {
		if err := i.encode(dst); err != nil {
			return err
		}
		if len(v.Declarations)-1 > idx {
			dst.WriteByte(semicolon)
		}
	}
	dst.WriteByte(rightCurlyBracket)

	return nil
}
