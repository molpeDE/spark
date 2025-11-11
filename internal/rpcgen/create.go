package rpcgen

import (
	"log"
	"net/http"
	"reflect"

	"github.com/molpeDE/spark/pkg/utils/validate"

	"github.com/fxamacker/cbor"
	"github.com/labstack/echo/v4"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
var ctxInterface = reflect.TypeOf((*echo.Context)(nil)).Elem()

func isValidHandler(meth reflect.Method) bool {
	return !(meth.Type.NumIn() < 2 ||
		meth.Type.NumIn() > 3 ||
		meth.Type.NumOut() != 2 ||
		!meth.Type.In(1).Implements(ctxInterface) ||
		!meth.Type.Out(1).Implements(errorInterface) ||
		meth.Type.Out(0).Kind() == reflect.Chan)
}

var allowedSyntax = "func(c echo.Context[, args Serializable]) (Serializable, error)"

type handlerInfo struct {
	method        reflect.Method
	typeofHandler reflect.Type
	hasArg        bool
	argType       reflect.Type
}

type bindHelper struct {
	obj          reflect.Value
	handlers     []*handlerInfo
	tsTypesToGen []reflect.Type
}

type RPC interface {
	Attach(*echo.Group)
	Typedefs() string
}

func From(rpcStruct any) RPC {

	var rpcHandler = reflect.TypeOf(rpcStruct)

	helper := &bindHelper{
		obj:          reflect.ValueOf(rpcStruct),
		handlers:     make([]*handlerInfo, 0, rpcHandler.NumMethod()),
		tsTypesToGen: make([]reflect.Type, 0),
	}

	for i := range rpcHandler.NumMethod() {
		meth := rpcHandler.Method(i)

		if !isValidHandler(meth) {
			log.Printf("WARNING: RPC Method '%s' is not in format %s", meth.Name, allowedSyntax)
			continue
		}

		hi := &handlerInfo{
			method:        meth,
			typeofHandler: rpcHandler,
			hasArg:        meth.Type.NumIn() == 3,
		}

		if hi.hasArg {
			hi.argType = meth.Type.In(2)
			helper.addIfStruct(hi.argType)
		}

		helper.addIfStruct(hi.method.Type.Out(0))

		helper.handlers = append(helper.handlers, hi)
	}

	return helper
}

func fail(c echo.Context, err error) error {
	c.Response().Header().Add("RPC-Failed", "1")
	return c.String(http.StatusBadRequest, err.Error())
}

func finalizer(c echo.Context, ret []reflect.Value) error {

	if !ret[1].IsNil() {
		return fail(c, ret[1].Interface().(error))
	}

	c.Response().Header().Add("content-type", "application/cbor")

	if err := cbor.NewEncoder(c.Response(), cbor.EncOptions{}).Encode(ret[0].Interface()); err != nil {
		return fail(c, err)
	}

	return nil
}

func (b *bindHelper) Attach(g *echo.Group) {
	for _, v := range b.handlers {

		if v.hasArg {
			g.POST("/"+v.method.Name, func(c echo.Context) error {
				args := make([]reflect.Value, 3)
				args[0] = b.obj
				args[1] = reflect.ValueOf(c)

				instance := reflect.New(v.argType)
				args[2] = instance.Elem()

				if err := cbor.NewDecoder(c.Request().Body).Decode(instance.Interface()); err != nil {
					return fail(c, err)
				}

				// TODO: check if needed
				// if instance.Elem().IsNil() {
				// 	return fail(c, fmt.Errorf("got nil when argument was expected"))
				// }

				if err := validate.ValidateStruct(instance.Elem()); err != nil {
					return fail(c, err)
				}

				return finalizer(c, v.method.Func.Call(args))
			})
		} else {
			g.POST("/"+v.method.Name, func(c echo.Context) error {
				return finalizer(c, v.method.Func.Call([]reflect.Value{b.obj, reflect.ValueOf(c)}))
			})
		}
	}
}
