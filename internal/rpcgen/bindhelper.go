//go:build !debug

package rpcgen

import "reflect"

func (*bindHelper) addIfStruct(reflect.Type) {}
func (*bindHelper) Typedefs() string         { return "" }
