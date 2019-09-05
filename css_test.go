package css

import (
	"reflect"
	"testing"
)

func TestDeclaration_String(t *testing.T) {
	type fields struct {
		Property string
		Value    []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "background",
			fields: fields{
				Property: "background",
				Value:    []string{"red"},
			},
			want: "background:red",
		},
		{
			name: "border",
			fields: fields{
				Property: "border",
				Value:    []string{"1px", "solid", "red"},
			},
			want: "border:1px solid red",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Declaration{
				Property: tt.fields.Property,
				Value:    tt.fields.Value,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("Declaration.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleset_String(t *testing.T) {
	type fields struct {
		Selectors    []Selector
		Declarations []Declaration
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Declarations: []Declaration{
					{
						Property: "background",
						Value:    []string{"red"},
					},
					{
						Property: "color",
						Value:    []string{"green"},
					},
				},
			},
			want: "background:red;color:green",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Ruleset{
				Selectors:    tt.fields.Selectors,
				Declarations: tt.fields.Declarations,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("Ruleset.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimple_Encode(t *testing.T) {
	type fields struct {
		Element        []byte
		Classes        [][]byte
		Attributes     []Attribute
		PseudoElements []Pseudo
		PseudoClasses  []Pseudo
		Negations      []Simple
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			fields: fields{
				Element: []byte("a"),
			},
			want: []byte("a"),
		},
		{
			fields: fields{
				Element: []byte("a"),
				Classes: [][]byte{[]byte("test")},
			},
			want: []byte("a.test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Simple{
				Element:        tt.fields.Element,
				Classes:        tt.fields.Classes,
				Attributes:     tt.fields.Attributes,
				PseudoElements: tt.fields.PseudoElements,
				PseudoClasses:  tt.fields.PseudoClasses,
				Negations:      tt.fields.Negations,
			}
			if got := v.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Simple.Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPseudo_Encode(t *testing.T) {
	type fields struct {
		Ident []byte
		Func  []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			fields: fields{
				Ident: []byte("nth-child"),
				Func:  []byte("4n"),
			},
			want: []byte("nth-child(4n)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Pseudo{
				Ident: tt.fields.Ident,
				Func:  tt.fields.Func,
			}
			if got := v.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pseudo.Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}
