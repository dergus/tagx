# TagX
Extract tag values from any Go struct fields

## Install
```sh
go get github.com/dergus/tagx
```

## Features

- Extract tag values by a tag key as a hierarchical  data structure
- Extract tag values by a tag key as a flat array of strings concatenated by a given delimiter

## Usage

### Extract
```go
type S struct {
    FieldA       int `json:"field_a"`
    NestedStruct struct {
    NestedFieldA        string `json:"nested_field_a"`
        AnotherNestedStruct struct {
        AnotherNestedField int `json:"another_nested_field"`
        } `json:"another_nested_struct"`
    } `json:"nested_struct"`
}

tags := tagx.Extract(S{}, "json")

// tags
[]tagx.Tag{
	{
		Value:"field_a",
	},
	{
		Value:"nested_struct",
		Children: []tagx.Tag{
			{
				Value: "nested_field_a",
			},
			{
				Value: "another_nested_struct",
				Children: []tagx.Tag{
					{
						Value:"another_nested_field",
					}
                }
            }
        }
	}
}
```

### ExtractFlat

```go
type S struct {
    FieldA       int `json:"field_a"`
    NestedStruct struct {
    NestedFieldA        string `json:"nested_field_a"`
        AnotherNestedStruct struct {
        AnotherNestedField int `json:"another_nested_field"`
        } `json:"another_nested_struct"`
    } `json:"nested_struct"`
}

tags := tagx.ExtractFlat(S{}, "json", ".")

// tags
[]string{
    "field_a",
    "nested_struct.nested_field_a",
    "nested_struct.another_nested_struct.another_nested_field",
}
```

## Caveats

- Only tag value is returned, options(e.g. `omitempty` for the json tag) are omitted
- Private Fields (with tag value `-`) are ignored and not returned
- Empty tag values or just whitespace values are ignored and not returned
- Fields of recursive structs aren't explored
## License
Copyright © 2023 [Ziyadin Shemsedinov](https://github.com/dergus).

This project is [MIT](./LICENSE) licensed.