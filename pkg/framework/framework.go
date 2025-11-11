package framework

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"unicode"

	"github.com/molpeDE/spark/internal/rpcgen"
	"github.com/molpeDE/spark/internal/spa"

	"github.com/labstack/echo/v4"
	"gopkg.in/ini.v1"
)

type MolpeApp interface {
	Attach(*echo.Group)
	SPA(fs.ReadDirFS) http.Handler
}

type molpeApp struct {
	instance any
}

func (m *molpeApp) Attach(g *echo.Group) {
	rpcgen.From(m.instance).Attach(g)
}

func (m *molpeApp) SPA(dist fs.ReadDirFS) http.Handler {
	return spa.SPA(dist, rpcgen.From(m.instance).Typedefs())
}

func CreateApp(rpcStruct any) MolpeApp {
	return &molpeApp{instance: rpcStruct}
}

/*
	Config Parser
*/

func canonicalName(fieldName string) string {
	var result []rune

	for i, r := range fieldName {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

func processStruct(tStruct reflect.Type, vStruct reflect.Value, cfgFile *ini.File, section string) error {
	for i := 0; i < tStruct.NumField(); i++ {
		field := tStruct.Field(i)
		value := vStruct.Field(i)

		keyName := canonicalName(field.Name)

		if field.Type.Kind() == reflect.Struct {
			if e := processStruct(field.Type, value, cfgFile, keyName); e != nil {
				return e
			}
			continue
		}

		val := cfgFile.Section(section).Key(keyName).MustString(field.Tag.Get("default"))

		switch field.Type.Kind() {
		case reflect.String:
			value.SetString(val)
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(val); err == nil {
				value.SetBool(boolVal)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
				value.SetUint(uintVal)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
				value.SetInt(intVal)
			}
		case reflect.Float32, reflect.Float64:
			if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
				value.SetFloat(floatVal)
			}
		default:
			log.Printf("[WARN: ConfigParser] unsupported type %s for key %s", field.Type.Kind(), field.Name)
		}
	}

	return nil
}

func ParseConfig(path string, configStruct any) error {

	if _, err := os.Stat(path); err != nil {
		if file, err := os.Create(path); err != nil {
			return err
		} else {
			file.Close()
		}
	}

	cfgFile, err := ini.LoadSources(ini.LoadOptions{}, path)

	if err != nil {
		return err
	}

	tConfStruct := reflect.TypeOf(configStruct)
	vConfStruct := reflect.ValueOf(configStruct)

	if tConfStruct.Kind() == reflect.Ptr {
		tConfStruct = tConfStruct.Elem()
		vConfStruct = vConfStruct.Elem()
	}

	if tConfStruct.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", tConfStruct.Kind())
	}

	if e := processStruct(tConfStruct, vConfStruct, cfgFile, ini.DefaultSection); e != nil {
		return e
	}

	if err := cfgFile.SaveToIndent(path, "\t"); err != nil {
		return err
	}

	return nil
}
