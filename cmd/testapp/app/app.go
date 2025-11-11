package app

import (
	"fmt"
	"log"
	"time"

	"github.com/molpeDE/spark/cmd/testapp/types"
	"github.com/molpeDE/spark/pkg/framework"

	"github.com/labstack/echo/v4"
)

type EchoResponse struct {
	Message string `json:"message"`
}

type TimeResponse struct {
	Time time.Time
}

type EchoRequest struct {
	Message string `json:"message" validate:"required"`
}

type App struct{}

var Instance = framework.CreateApp(&App{})

func (*App) Example(c echo.Context, req EchoRequest) (EchoResponse, error) {
	return EchoResponse{Message: fmt.Sprintf("echo: %s", req.Message)}, nil
}

func (*App) ExamplePlain(c echo.Context, req EchoRequest) (string, error) {
	return fmt.Sprintf("echo: %s", req.Message), nil
}

func (*App) FailForMe(c echo.Context, req EchoRequest) (string, error) {
	return "", fmt.Errorf("error: %s", req.Message)
}

func (*App) BinaryExample(c echo.Context) ([]byte, error) {
	return []byte{0x1, 0x2}, nil
}

func (*App) NativeTypeExample(c echo.Context, in []float32) (float64, error) {
	for _, v := range in {
		log.Println(v)
	}
	return 0.4444, nil
}

func (*App) GetTime(c echo.Context) (TimeResponse, error) {
	return TimeResponse{Time: time.Now()}, nil
}

func (*App) RandBytes(c echo.Context) ([]byte, error) {
	return []byte{0x1, 0x1, 0x1}, nil
}

func (*App) TypeHandling(c echo.Context) (types.ExtendedField, error) {
	return types.ExtendedField{}, nil
}

func (*App) TypeHandling2(c echo.Context) (types.TestStruct, error) {
	return types.TestStruct{}, nil
}
