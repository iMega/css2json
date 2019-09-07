package css2json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	js := []byte(`[
		{
	    	"atrule":{
				"ident":{
					"type":"charset",
					"info": {"value":"asf"}
				}
			},
		  "ruleset": {
			"selectors": [
			  {
				"simple": {
				  "element": "p"
				}
			  }
			],
			"declarations": [
			  {
				"property": "color",
				"value": [
				  "red"
				]
			  },
			  {
				"property": "border",
				"value": [
				  "1px",
				  "solid",
				  "red"
				]
			  }
			]
		  }
		}
	  ]`)

	a := Statements{}
	json.Unmarshal(js, &a)
	buf := &bytes.Buffer{}
	a[0].Ruleset.encode(buf)
	got := buf.String()

	want := `p{color:red;border:1px solid red}`
	if want != got {
		t.Errorf("TestUnmarshalJSON() = %v, want %v", got, want)
	}
}

func TestMarshalJSON(t *testing.T) {
	s := Statements{
		{
			AtRule: &AtRule{
				Identifier: Identifier{
					Type: TextBytes("charset"),
					Information: &CharsetInformation{
						Value: TextBytes("utf-8"),
					},
				},
			},
			Ruleset: &Ruleset{
				Selectors: []Selector{
					{
						Simple: Simple{
							Element: []byte("p"),
						},
					},
				},
				Declarations: []Declaration{
					{
						Property: TextBytes("color"),
						Value: []TextBytes{
							TextBytes("red"),
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(s)
	if err != nil {
		t.Errorf("TestMarshalJSON error = %v", err)
	}

	want := `[{"atrule":{"ident":{"type":"charset","info":{"value":"utf-8"}}},"ruleset":{"selectors":[{"simple":{"element":"p"}}],"declarations":[{"property":"color","value":["red"]}]}}]`
	got := string(b)
	if got != want {
		t.Errorf("TestMarshalJSON() = %v, want %v", got, want)
	}
}

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
		Element        TextBytes
		Classes        []TextBytes
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
		{
			fields: fields{
				Element: []byte("a"),
				Classes: []TextBytes{
					[]byte("myclass"),
				},
				Attributes: []Attribute{
					{
						Attr:     []byte("href"),
						Operator: []byte("*="),
						Value:    []byte(".com"),
						Modifier: []byte("s"),
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `a.myclass[href*=".com" s]`,
			wantErr: false,
		},
		{
			fields: fields{
				Element: []byte("a"),
				Classes: []TextBytes{
					TextBytes("myclass"),
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `a.myclass`,
			wantErr: false,
		},
		{
			fields: fields{
				Element: []byte("html|*"),
				Negations: []Simple{
					{
						PseudoClasses: []Pseudo{
							{
								Ident: []byte("link"),
							},
						},
					},
					{
						PseudoClasses: []Pseudo{
							{
								Ident: []byte("visited"),
							},
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `html|*:not(:link):not(:visited)`,
			wantErr: false,
		},
		{
			fields: fields{
				Element: []byte("button"),
				Negations: []Simple{
					{
						Attributes: []Attribute{
							{
								Attr: []byte("DISABLED"),
							},
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `button:not([DISABLED])`,
			wantErr: false,
		},
		{
			fields: fields{
				Element: []byte("*"),
				Negations: []Simple{
					{
						Element: []byte("FOO"),
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `*:not(FOO)`,
			wantErr: false,
		},
		{
			fields: fields{
				Element: []byte("tr"),
				PseudoClasses: []Pseudo{
					{
						Ident: []byte("nth-child"),
						Func:  []byte("2n+1"),
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `tr:nth-child(2n+1)`,
			wantErr: false,
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
		Property TextBytes
		Value    []TextBytes
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
				Property: TextBytes("color"),
				Value: []TextBytes{
					TextBytes("red"),
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `color:red`,
			wantErr: false,
		},
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
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Declaration.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelector_encode(t *testing.T) {
	type fields struct {
		Simple     Simple
		Combinates []Combinate
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
				Simple: Simple{
					Element: []byte("div"),
				},
			},
			want: `div`,
			args: args{
				dst: &bytes.Buffer{},
			},
			wantErr: false,
		},
		{
			fields: fields{
				Simple: Simple{
					Element: []byte("div"),
				},
				Combinates: []Combinate{
					{
						Combinator: []byte(">"),
						Simple: Simple{
							Element: []byte("p"),
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `div>p`,
			wantErr: false,
		},
		{
			fields: fields{
				Simple: Simple{
					Element: []byte("div"),
				},
				Combinates: []Combinate{
					{
						Combinator: []byte(">"),
						Simple: Simple{
							Element: []byte("p"),
						},
					},
					{
						Combinator: []byte("~"),
						Simple: Simple{
							Element: []byte("a"),
						},
					},
					{
						Combinator: []byte(" "),
						Simple: Simple{
							Element: []byte("span"),
						},
					},
					{
						Combinator: []byte("+"),
						Simple: Simple{
							Element: []byte("b"),
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `div>p~a span+b`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Selector{
				Simple:     tt.fields.Simple,
				Combinates: tt.fields.Combinates,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Selector.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Selector.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCombinate_encode(t *testing.T) {
	type fields struct {
		Combinator []byte
		Simple     Simple
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
				Combinator: []byte(">"),
				Simple: Simple{
					Element: []byte("p"),
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `>p`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Combinate{
				Combinator: tt.fields.Combinator,
				Simple:     tt.fields.Simple,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Combinate.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combinate.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleset_encode(t *testing.T) {
	type fields struct {
		Selectors    []Selector
		Declarations []Declaration
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
				Selectors: []Selector{
					{
						Simple: Simple{
							Element: []byte("p"),
						},
					},
					{
						Simple: Simple{
							Element: []byte("span"),
						},
					},
				},
				Declarations: []Declaration{
					{
						Property: []byte("color"),
						Value: []TextBytes{
							TextBytes("red"),
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `p,span{color:red}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Ruleset{
				Selectors:    tt.fields.Selectors,
				Declarations: tt.fields.Declarations,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Ruleset.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ruleset.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextBytes_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		v       TextBytes
		want    []byte
		wantErr bool
	}{
		{
			v:       TextBytes("red"),
			want:    []byte(`"red"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TextBytes.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextBytes.MarshalJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		s Statements
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				s: Statements{
					{
						AtRule: &AtRule{
							Identifier: Identifier{
								Type: TextBytes("charset"),
								Information: &CharsetInformation{
									Value: TextBytes("utf-8"),
								},
							},
						},
					},
					{
						Ruleset: &Ruleset{
							Selectors: []Selector{
								{
									Simple: Simple{
										Element: TextBytes("span"),
									},
								},
							},
							Declarations: []Declaration{
								{
									Property: TextBytes("color"),
									Value: []TextBytes{
										TextBytes("red"),
									},
								},
							},
						},
					},
				},
			},
			want:    `@charset "utf-8";span{color:red};`,
			wantErr: false,
		},
		{
			args: args{
				s: Statements{
					{
						AtRule: &AtRule{
							Identifier: Identifier{
								Type: TextBytes("media"),
								Information: &MediaInformation{
									Queries: []Query{
										{
											Type: &Type{
												Value: TextBytes("all"),
											},
											Conditions: []Condition{
												{
													Operator: TextBytes("and"),
													Feature:  TextBytes("max-width"),
													Value:    TextBytes("699px"),
												},
												{
													Operator: TextBytes("and"),
													Feature:  TextBytes("max-width"),
													Value:    TextBytes("520px"),
												},
											},
										},
										{
											Conditions: []Condition{
												{
													Feature: TextBytes("max-width"),
													Value:   TextBytes("1151px"),
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Ruleset: &Ruleset{
							Selectors: []Selector{
								{
									Simple: Simple{
										Element: TextBytes("#sidebar"),
									},
									Combinates: []Combinate{
										{
											Combinator: TextBytes(" "),
											Simple: Simple{
												Element: TextBytes("ul"),
											},
										},
										{
											Combinator: TextBytes(" "),
											Simple: Simple{
												Element: TextBytes("li"),
											},
										},
										{
											Combinator: TextBytes(" "),
											Simple: Simple{
												Element: TextBytes("a"),
											},
										},
									},
								},
							},
							Declarations: []Declaration{
								{
									Property: TextBytes("padding-left"),
									Value: []TextBytes{
										TextBytes("21px"),
									},
								},
								{
									Property: TextBytes("background"),
									Value: []TextBytes{
										TextBytes("url(../images/email.png)"),
										TextBytes("left"),
										TextBytes("center"),
										TextBytes("no-repeat"),
									},
								},
							},
						},
					},
				},
			},
			want:    `@media all and (max-width: 699px) and (min-width: 520px),(min-width: 1151px){#sidebar ul li a{padding-left:21px;background:url(../images/email.png) left center no-repeat}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStatement_encode(t *testing.T) {
	type fields struct {
		AtRule  *AtRule
		Ruleset *Ruleset
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
				AtRule: &AtRule{
					Identifier: Identifier{
						Type: TextBytes("charset"),
						Information: &CharsetInformation{
							Value: TextBytes("utf-8"),
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `@charset "utf-8"`,
		},
		{
			fields: fields{
				Ruleset: &Ruleset{
					Selectors: []Selector{
						{
							Simple: Simple{
								Element: []byte("p"),
							},
						},
						{
							Simple: Simple{
								Element: []byte("span"),
							},
						},
					},
					Declarations: []Declaration{
						{
							Property: []byte("color"),
							Value: []TextBytes{
								TextBytes("red"),
							},
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `p,span{color:red}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Statement{
				AtRule:  tt.fields.AtRule,
				Ruleset: tt.fields.Ruleset,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Statement.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Statement.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtRule_encode(t *testing.T) {
	type fields struct {
		Identifier Identifier
		Nested     []*Statement
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
				Identifier: Identifier{
					Type: TextBytes("charset"),
					Information: &CharsetInformation{
						Value: TextBytes("utf-8"),
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `@charset "utf-8"`,
		},
		{
			fields: fields{
				Identifier: Identifier{
					Type: TextBytes("keyframes"),
					Information: &KeyframesInformation{
						Value: TextBytes("slide-right"),
					},
				},
				Nested: []*Statement{
					{
						Ruleset: &Ruleset{
							Selectors: []Selector{
								{
									Simple: Simple{
										Element: TextBytes("from"),
									},
								},
							},
							Declarations: []Declaration{
								{
									Property: TextBytes("margin-left"),
									Value: []TextBytes{
										TextBytes("0px"),
									},
								},
							},
						},
					},
					{
						Ruleset: &Ruleset{
							Selectors: []Selector{
								{
									Simple: Simple{
										Element: TextBytes("50%"),
									},
								},
							},
							Declarations: []Declaration{
								{
									Property: TextBytes("margin-left"),
									Value: []TextBytes{
										TextBytes("110px"),
									},
								},
								{
									Property: TextBytes("opacity"),
									Value: []TextBytes{
										TextBytes("0.9"),
									},
								},
							},
						},
					},
					{
						Ruleset: &Ruleset{
							Selectors: []Selector{
								{
									Simple: Simple{
										Element: TextBytes("to"),
									},
								},
							},
							Declarations: []Declaration{
								{
									Property: TextBytes("margin-left"),
									Value: []TextBytes{
										TextBytes("200px"),
									},
								},
							},
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `@keyframes slide-right{from{margin-left:0px}50%{margin-left:110px;opacity:0.9}to{margin-left:200px}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &AtRule{
				Identifier: tt.fields.Identifier,
				Nested:     tt.fields.Nested,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("AtRule.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AtRule.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharsetInformation_encode(t *testing.T) {
	type fields struct {
		Value TextBytes
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
				Value: TextBytes("utf-8"),
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want: `"utf-8"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &CharsetInformation{
				Value: tt.fields.Value,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CharsetInformation.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CharsetInformation.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifier_encode(t *testing.T) {
	type fields struct {
		Type        TextBytes
		Information Information
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
				Type: TextBytes("charset"),
				Information: &CharsetInformation{
					Value: TextBytes("utf-8"),
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `@charset "utf-8"`,
			wantErr: false,
		},
		{
			fields: fields{
				Type: TextBytes(""),
				Information: &CharsetInformation{
					Value: TextBytes("utf-8"),
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Identifier{
				Type:        tt.fields.Type,
				Information: tt.fields.Information,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Identifier.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Identifier.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyframesInformation_encode(t *testing.T) {
	type fields struct {
		Value TextBytes
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
				Value: TextBytes("name"),
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `name`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &KeyframesInformation{
				Value: tt.fields.Value,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("KeyframesInformation.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyframesInformation.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaInformation_encode(t *testing.T) {
	type fields struct {
		Queries []Query
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
				Queries: []Query{
					{
						Type: &Type{
							Value: TextBytes("screen"),
						},
						Conditions: []Condition{
							{
								Operator: TextBytes("and"),
								Feature:  TextBytes("max-width"),
								Value:    TextBytes("650px"),
							},
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `screen and (max-width:650px)`,
			wantErr: false,
		},
		{
			fields: fields{
				Queries: []Query{
					{
						Type: &Type{
							Value: TextBytes("screen"),
						},
					},
					{
						Type: &Type{
							Value: TextBytes("print"),
						},
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `screen,print`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &MediaInformation{
				Queries: tt.fields.Queries,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("MediaInformation.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MediaInformation.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_encode(t *testing.T) {
	type fields struct {
		Type       *Type
		Conditions []Condition
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
				Type: &Type{
					Value: TextBytes("screen"),
				},
				Conditions: []Condition{
					{
						Operator: TextBytes("and"),
						Feature:  TextBytes("max-width"),
						Value:    TextBytes("650px"),
					},
				},
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `screen and (max-width:650px)`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Query{
				Type:       tt.fields.Type,
				Conditions: tt.fields.Conditions,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Query.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_encode(t *testing.T) {
	type fields struct {
		Operator TextBytes
		Value    TextBytes
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
				Operator: TextBytes("only"),
				Value:    TextBytes("screen"),
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `only screen`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Type{
				Operator: tt.fields.Operator,
				Value:    tt.fields.Value,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Type.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCondition_encode(t *testing.T) {
	type fields struct {
		Operator TextBytes
		Feature  TextBytes
		Value    TextBytes
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
				Operator: TextBytes("and"),
				Feature:  TextBytes("max-width"),
				Value:    TextBytes("650px"),
			},
			args: args{
				dst: &bytes.Buffer{},
			},
			want:    `and (max-width:650px)`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Condition{
				Operator: tt.fields.Operator,
				Feature:  tt.fields.Feature,
				Value:    tt.fields.Value,
			}
			if err := v.encode(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Condition.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.dst.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Condition.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
