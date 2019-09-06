package css

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPseudo_encode(t *testing.T) {
	type fields struct {
		Ident []byte
		Func  []byte
	}
	type args struct {
		dst *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				Ident: []byte("nth-child"),
				Func:  []byte("4n"),
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    "nth-child(4n)",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Pseudo{
				Ident: tt.fields.Ident,
				Func:  tt.fields.Func,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Pseudo.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pseudo.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttribute_encode(t *testing.T) {
	type fields struct {
		Attr     []byte
		Operator []byte
		Value    []byte
		Modifier []byte
	}
	type args struct {
		dst *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				Attr:     []byte("a"),
				Operator: []byte("*="),
				Value:    []byte(".com"),
				Modifier: []byte("s"),
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `[a*=".com" s]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Attribute{
				Attr:     tt.fields.Attr,
				Operator: tt.fields.Operator,
				Value:    tt.fields.Value,
				Modifier: tt.fields.Modifier,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Attribute.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Attribute.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimple_encode(t *testing.T) {
	type fields struct {
		Element        []byte
		Classes        [][]byte
		Attributes     []Attribute
		PseudoElements []Pseudo
		PseudoClasses  []Pseudo
		Negations      []Simple
	}
	type args struct {
		dst *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Simple.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Simple.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeclaration_encode(t *testing.T) {
	type fields struct {
		Property []byte
		Value    [][]byte
	}
	type args struct {
		dst *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Declaration{
				Property: tt.fields.Property,
				Value:    tt.fields.Value,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Declaration.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
