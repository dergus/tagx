package tagx

import (
	"reflect"
	"strings"
)

type Tag struct {
	Value    string
	Children []Tag
}

const privateFieldVal = "-"

func Extract(structure any, tagKey string) []Tag {
	return extract(reflect.TypeOf(structure), tagKey, nil)
}

func extract(typ reflect.Type, tagKey string, seen func(p reflect.Type) bool) []Tag {
	var tags []Tag

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	// avoid recursive data structures
	if seen != nil && seen(typ) {
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		rawTag := field.Tag.Get(tagKey)
		tagParts := strings.Split(rawTag, ",")

		if len(tagParts) == 0 {
			continue
		}

		if tagParts[0] == privateFieldVal {
			continue
		}

		if strings.TrimSpace(tagParts[0]) == "" {
			continue
		}

		tags = append(tags, Tag{
			Value: tagParts[0],
			Children: extract(field.Type, tagKey, func(childType reflect.Type) bool {
				if childType.PkgPath() == "" || childType.Name() == "" {
					return false
				}

				if childType.PkgPath() == typ.PkgPath() && childType.Name() == typ.Name() {
					return true
				}

				return seen != nil && seen(childType)
			}),
		})
	}

	return tags
}

func ExtractFlat(structure any, tagKey, delimiter string) []string {
	tags := Extract(structure, tagKey)

	return extractFlat(tags, delimiter)
}

func extractFlat(tags []Tag, delimiter string) []string {
	var tagVals []string

	for _, t := range tags {
		if len(t.Children) == 0 {
			tagVals = append(tagVals, t.Value)
			continue
		}

		for _, v := range extractFlat(t.Children, delimiter) {
			tagVals = append(tagVals, t.Value+delimiter+v)
		}
	}

	return tagVals
}
