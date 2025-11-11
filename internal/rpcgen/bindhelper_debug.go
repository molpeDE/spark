//go:build debug

package rpcgen

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/molpeDE/spark/internal/tsgen"
)

func (b *bindHelper) addIfStruct(t reflect.Type) {
	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		b.tsTypesToGen = append(b.tsTypesToGen, t)
	} else if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct {
		b.tsTypesToGen = append(b.tsTypesToGen, t.Elem())
	} else if t.Kind() == reflect.Struct {
		b.tsTypesToGen = append(b.tsTypesToGen, t)
	}
}

func mkJSType(t reflect.Type) string {
	kind := t.Kind()
	switch kind {
	case reflect.Ptr:
		return mkJSType(t.Elem())
	case reflect.Bool:
		return "boolean"
	case reflect.Slice:
		jsType := fmt.Sprintf("%s[]", mkJSType(t.Elem()))
		if jsType == "number /* uint8 */[]" {
			return "Uint8Array" // special case
		}
		return jsType
	case reflect.String:
		return "string"
	case reflect.Struct, reflect.Interface:
		return t.Name()
	}

	switch t.String() {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"complex64", "complex128",
		"rune":
		return "number /* " + t.String() + " */"
	}

	return "any"
}

func (b *bindHelper) Typedefs() string {
	var builder strings.Builder

	builder.WriteString("// AUTOMATICALLY GENERATED - DO NOT EDIT\n")

	var converter = tsgen.New()
	converter.CreateInterface = true
	converter.CreateConstructor = false
	converter.CreateFromMethod = false
	converter.Indent = "\t"
	converter.ManageType(time.Time{}, tsgen.TypeOptions{TSType: "number /* unix timestamp */"})
	converter.ManageType([]byte{}, tsgen.TypeOptions{TSType: "Uint8Array"})

	for _, v := range b.tsTypesToGen {
		converter.AddType(v)
	}

	defs, err := converter.Convert(map[string]string{})

	if err != nil {
		log.Printf("Warn: %s", err)
		return ""
	}

	builder.WriteString(defs)

	builder.WriteString(fmt.Sprintf("\n\n/// RPC Generated\nexport interface %s {\n", b.obj.Type().Elem().Name()))

	for _, hi := range b.handlers {
		builder.WriteString(fmt.Sprintf("\t%s(", hi.method.Name))

		if hi.hasArg {
			builder.WriteString(fmt.Sprintf("arg0: %s", mkJSType(hi.argType)))
		}

		builder.WriteString(fmt.Sprintf("): Promise<%s>\n", mkJSType(hi.method.Type.Out(0))))
	}

	builder.WriteString("}")

	return builder.String()
}
