package helper

import (
	"fmt"
	"reflect"
)

func GeStructJsonFields(structs ...interface{}) []string {
	fields := []string{}
	seenFields := make(map[string]struct{}) // To store unique fields

	for _, s := range structs {
		t := reflect.TypeOf(s)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			fmt.Printf("Warning: Skipping non-struct type %T\n", s)
			continue
		}

		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("json")
			if tag != "" && tag != "-" {
				if _, found := seenFields[tag]; !found {
					fields = append(fields, tag)
					seenFields[tag] = struct{}{}
				}
			}
		}
	}
	return fields
}
