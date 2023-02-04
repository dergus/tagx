package tagx

import (
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	type args struct {
		structure any
		tagKey    string
	}

	type innerStruct struct {
		NestedFieldA int `json:"nested_field_a"`
		NestedFieldB int `json:"nested_field_b"`
	}

	type nestedStruct struct {
		JustField     string      `json:"just_field"`
		AnotherStruct innerStruct `json:"another_struct"`
	}

	type testStruct struct {
		SimpleField            string       `json:"simple_field"`
		SimpleFieldWithOptions string       `json:"simple_field_with_options,omitempty"`
		NestedStruct           innerStruct  `json:"nested_struct"`
		NestedPointerStruct    *innerStruct `json:"nested_pointer_struct"`
		DeepStruct             nestedStruct `json:"deep_struct"`
		AnonNestedStruct       struct {
			JustField               int `json:"just_field"`
			AnotherAnonNestedStruct struct {
				JustField int `json:"just_field"`
			} `json:"another_anon_nested_struct"`
		} `json:"anon_nested_struct"`
		SelfReference       *testStruct `json:"self_reference"`
		NestedSelfReference struct {
			JustField string      `json:"just_field"`
			SelfRef   *testStruct `json:"self_ref"`
		} `json:"nested_self_reference"`
		FieldWithoutTags             float64
		FieldWithoutNeededTag        float64 `xml:"field_without_needed_tag"`
		FieldWithMultipleTags        bool    `xml:"field_with_multiple_tags" json:"field_with_multiple_tags"`
		FieldWithEmptyTagValue       int     `json:""`
		FieldWithCommaInsteadOfValue int     `json:","`
		FieldWithWhiteSpaceValue     string  `json:"    "`
		FieldWithPrivateValue        int     `json:"-"`
		FieldWithJustOptions         int     `json:",omitempty"`
		FieldWithInvalidTag          int     `json:field_with_invalid_tag`
	}

	expected := []Tag{
		{Value: "simple_field"},
		{Value: "simple_field_with_options"},
		{
			Value: "nested_struct",
			Children: []Tag{
				{Value: "nested_field_a"},
				{Value: "nested_field_b"},
			},
		},
		{
			Value: "nested_pointer_struct",
			Children: []Tag{
				{Value: "nested_field_a"},
				{Value: "nested_field_b"},
			},
		},
		{
			Value: "deep_struct",
			Children: []Tag{
				{Value: "just_field"},
				{
					Value: "another_struct",
					Children: []Tag{
						{Value: "nested_field_a"},
						{Value: "nested_field_b"},
					},
				},
			},
		},
		{
			Value: "anon_nested_struct",
			Children: []Tag{
				{Value: "just_field"},
				{
					Value: "another_anon_nested_struct",
					Children: []Tag{
						{Value: "just_field"},
					},
				},
			},
		},
		{Value: "self_reference"},
		{
			Value: "nested_self_reference",
			Children: []Tag{
				{Value: "just_field"},
				{Value: "self_ref"},
			},
		},
		{Value: "field_with_multiple_tags"},
	}

	tests := []struct {
		name string
		args args
		want []Tag
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    "json",
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    "json",
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Extract(tt.args.structure, tt.args.tagKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractFlat(t *testing.T) {
	type args struct {
		structure any
		tagKey    string
		delimiter string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should return flat array of tag values",
			args: args{
				structure: struct {
					FieldA       int `json:"field_a"`
					NestedStruct struct {
						NestedFieldA        string `json:"nested_field_a"`
						AnotherNestedStruct struct {
							AnotherNestedField int `json:"another_nested_field"`
						} `json:"another_nested_struct"`
					} `json:"nested_struct"`
				}{},
				tagKey:    "json",
				delimiter: ".",
			},
			want: []string{
				"field_a",
				"nested_struct.nested_field_a",
				"nested_struct.another_nested_struct.another_nested_field",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractFlat(tt.args.structure, tt.args.tagKey, tt.args.delimiter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractFlat() = %v, want %v", got, tt.want)
			}
		})
	}
}
