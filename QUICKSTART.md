# Spark — Quickstart

This quickstart shows the minimal steps to use the Spark framework in a new blank Go application. It covers:

- creating a simple RPC service,
- parsing configuration with defaults,
- attaching RPC routes to an Echo server,
- serving a SPA (production with embedded assets and development mode with live bundler),
- calling RPC methods from a TypeScript frontend using the built-in client.

Prerequisites
- Go (1.18+)
- bun (https://bun.sh) installed and in PATH

Project skeleton
```
myapp/
├─ go.mod
├─ cmd/
│  └─ myapp/
│     └─ main.go
├─ pkg/
│  └─ app/
│     └─ app.go
├─ cfg/
│  └─ cfg.go
└─ frontend/
   ├─ index.html
   ├─ app.tsx
   └─ (dist or dev files)
```

1) Create module and add dependency
```bash
mkdir myapp
cd myapp
go mod init example.com/myapp

# add spark (replace with the canonical module path if different)
go get github.com/molpeDE/spark@latest
```

2) Minimal RPC app implementation

Create `pkg/app/app.go`. RPC methods must use one of the supported signatures:
- func (a *App) Method(c echo.Context) (SomeSerializable, error)
- func (a *App) Method(c echo.Context, arg SomeSerializable) (SomeSerializable, error)

The framework will automatically register methods that match the expected signature and will decode/encode request/response using CBOR. Validation tags are supported on argument structs via `go-playground/validator`.

Example:
```go
package app

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/molpeDE/spark/pkg/framework"
)

type TimeResp struct {
	Time int64 `json:"time"`
}

type EchoReq struct {
	Message string `json:"message" validate:"required"`
}

type App struct{}

var Instance = framework.CreateApp(&App{})

func (*App) GetTime(c echo.Context) (TimeResp, error) {
	return TimeResp{Time: time.Now().Unix()}, nil
}

func (*App) Echo(c echo.Context, req EchoReq) (string, error) {
	return fmt.Sprintf("echo: %s", req.Message), nil
}
```

3) Configuration with defaults

Create `cfg/cfg.go`. Spark provides ParseConfig(path, structPtr) which reads/creates an INI and fills struct fields. Use struct tags `default:"..."`.

Example:
```go
package cfg

import "github.com/molpeDE/spark/pkg/framework"

type Config struct {
	Server struct {
		Host string `default:"0.0.0.0"`
		Port uint16 `default:"8080"`
	}
	Production bool    `default:"false"`
	SomeFloat  float64 `default:"1.0"`
}

var inst = &Config{}

func Parse(path string) error {
	return framework.ParseConfig(path, inst)
}

func Get() Config {
	return *inst
}
```

When you first run the app, the file is created if missing and default values are written.

4) Main — hook everything up and serve the SPA

Create `cmd/myapp/main.go`. Attach the RPC group and a handler for serving frontend assets. You can either embed static `dist` files (production) or use the development SPA proxy (debug build tag) which starts the frontend bundler and proxies to it.

Minimal main:
```go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"example.com/myapp/pkg/app"
	"example.com/myapp/cfg"
	// import your frontend package which should expose a variable SPA http.Handler
	"example.com/myapp/frontend"
)

func main() {
	configPath := flag.String("config", "./config.ini", "config file path")
	flag.Parse()

	if err := cfg.Parse(*configPath); err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// serve frontend (either embedded files or dev server depending on build tags)
	e.GET("/*", echo.WrapHandler(frontend.SPA), middleware.Gzip())

	// attach RPC endpoints under /rpc
	backend := e.Group("/rpc")
	app.Instance.Attach(backend)

	addr := fmt.Sprintf("%s:%d", cfg.Get().Server.Host, cfg.Get().Server.Port)
	e.Logger.Fatal(e.Start(addr))
}
```

5) Frontend — production embedding and dev mode

Production: build your SPA into `dist/` and embed it into your Go binary. Example `frontend` package that embeds a `dist` dir:

frontend/frontend.go:
```go
package frontend

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/molpeDE/spark/cmd/testapp/app"
)

//go:generate env NODE_PATH=node_modules bun node_modules/@molpe/spark/bundler --prod
//go:embed all:dist
var _distEmbed embed.FS

func init() {
	var subfs, _ = fs.Sub(_distEmbed, "dist")
	SPA = app.Instance.SPA(subfs.(fs.ReadDirFS))
}

var SPA http.Handler
```

(See the framework's `internal/spa/embedded_assets.go` for the expected behaviour. The test app in the repo shows a complete example; you can copy that pattern.)

Dev mode: Build with the debug build tag and the framework will start a bundler/proxy during development which picks up TypeScript types (automatic generation) and hot-reload.

Run in dev:
```bash
go run -tags=debug ./cmd/myapp
```

Note: The debug SPA backend expects `bun` and the internal bundler to be available. The repo's testapp demonstrates the full dev workflow.

6) Calling RPC from the frontend (TypeScript)

Spark ships a small client helper in TypeScript that encodes requests as CBOR and decodes the response.

Example usage (client from @molpe/spark):
```ts
import Client from '@molpe/spark/client'
import type { App } from './gotypes' // generated types

const client = Client<App>('/rpc/')

async function call() {
  const resp = await client.Echo({ message: 'hello' })
  console.log(resp) // "echo: hello"
}

// in Preact/React you can also use the generated SWR hooks:
// const { data, isLoading } = client.useGetTime(undefined, { refreshInterval: 1000 })
```

7) Build & run

Production build (assumes you have built your frontend into `dist` and embedded it):
```bash
go generate ./...
go build -o myapp ./cmd/myapp
./myapp -config config.ini
```

Development (with live bundling / type generation):
```bash
# run with debug build tag; this will start the bundler proxy and generate gotypes.ts
go run -tags=debug ./cmd/myapp -config config.ini
```

8) Notes & tips
- RPC function signatures must match the allowed patterns; otherwise they are ignored and a warning is logged.
- Request and responses are encoded with CBOR. For JSON interchange you will need to write a wrapper, but the framework expects cbor for RPC calls.
- Validation: argument structs support `validate` tags. The framework will run validation and return a 400-like error with message on failure.
- Type generation: in debug builds the framework writes TypeScript typedefs (so your frontend can import `gotypes.ts`).
- Inspect the repository's `cmd/testapp` for a full, working example (app, config, frontend embed and dev setup).

That's it — you now have a minimal Spark application with RPC endpoints, config parsing, frontend serving, and TypeScript client integration. For more advanced use (JWT, middleware, custom validators, binary returns, arrays and complex types), inspect the example in the repository `cmd/testapp` which demonstrates many of the features.
