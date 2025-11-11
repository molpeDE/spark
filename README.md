# Spark

Spark is a small, batteries-included Go framework for building microservices with:
- Simple RPC-by-reflection
- Automatic TypeScript type generation for your RPC types
- Embedded Single-Page Application (SPA) support for shipping frontend assets with your binary
- Config parsing into Go structs with defaults
- Request validation using go-playground/validator
- Development bundler / hot-reload support for frontend development (powered by Bun)

This repository contains a demo application under `cmd/testapp` that exercises the framework features.

### Summary
- Reflection-based RPC exposure using Echo. Methods follow this signature:  
  `func (s *YourStruct) Method(c echo.Context[, args Serializable]) (Serializable, error)`
- Client: lightweight TypeScript frontend usign preact and generated RPC types
- Frontend bundling and hot-reload via Bun
- Config parsing from INI file into Go structs with defaults
- Validation using go-playground/validator

Contents
- cmd/testapp — example application demonstrating features
- pkg/framework — the core functionality of this package is exposed here
- pkg/utils/validate — helper for struct validation

Requirements
- Go 1.20+ (or latest supported)
- Bun (for frontend bundling and development server) if you want to run frontend dev mode
- Node tooling only if you use any npm/bun packages in the frontend, but Bun simplifies this
- git

### Quick start  
- Use in a project: [Click Me](./QUICKSTART.md)

- run example app (development)
   1. Clone the repository
      git clone https://github.com/molpeDE/spark.git
      cd spark

   2. Ensure Bun is installed (dev frontend only)
      - Install Bun from https://bun.sh (The debug/dev frontend expects Bun for the bundler and dev server.)

   3. Run the example app (uses cmd/testapp)
      go run ./cmd/testapp

      Notes:
      - If `config.ini` does not exist it will be created with any default values found in the app's config struct.
      - The testapp serves the SPA root (frontend) and exposes RPC endpoints under `/rpc`.

   4. Open the app in your browser (defaults from example config):
      http://localhost:3999/

- build production
   1. Build the frontend for production (project uses a Bun-based bundler):
      `go generate cmd/testapp/frontend/frontend.go`

   2. Build and run the test-app:  
      ```
      go build -o spark-example-app ./cmd/testapp
      ./spark-example-app -config ./config.ini
      ```

#### Additional information

Configuration
- Configuration parsing lives in pkg/framework and parses an INI file into a Go struct, honoring `default` tags on struct fields.
- Example config struct location: cmd/testapp/cfg/cfg.go
- ParseConfig will create the file if it doesn't exist and will populate it with defaults.

RPC (server-side)
- Define methods on a struct of the form:
  func (s *YourStruct) Method(c echo.Context, arg YourArg) (YourReturn, error)
  or
  func (s *YourStruct) Method(c echo.Context) (YourReturn, error)
- Methods are discovered via reflection (internal/rpcgen) and attached to an Echo group (see framework.CreateApp and rpcgen.From).
- Request/response encoding: CBOR
  - Content-Type must be "application/cbor" when calling RPC endpoints.
  - On error, responses include header `RPC-Failed: 1` and a plain text error message body.
- Error handling: RPC handlers should return an error as the second return value. The framework will convert that to a failed HTTP response.

RPC (client-side)
- A small TypeScript client is included at internal/typescript/client.ts. It:
  - Uses fetch to POST CBOR-encoded bytes to `/rpc/<MethodName>`.
  - Throws on RPC failures (when server sets `Rpc-Failed: 1` header).
  - Exposes generated "use" hooks in TypeScript when using the provided typedef generation in debug builds.
- In debug builds the Go code can generate TypeScript typedefs for RPC types (see internal/rpcgen/bindhelper_debug.go and internal/tsgen).

Frontend / SPA
- In debug builds the server will write generated types (gotypes.ts) into the frontend directory and spawn a Bun-based bundler + dev server (hot-reload).
- In non-debug builds the SPA uses embedded dist assets (go:embed) and serves static files, falling back to index.html for client-side routing.

Validation
- Request argument structs are validated using go-playground/validator. Add `validate` tags to struct fields where needed.
- Validation errors are returned as a textual error response.

WebSockets
- The project contains a small TypeScript helper for async WebSocket handling (internal/typescript/websocket.ts) used by the example frontend.

#### Development notes
- Debug builds enable additional behavior:
  - Generation of TypeScript typedefs from Go types (for a nicer developer DX).
  - Bun-based dev bundler runs and the server reverse-proxies to it.
- Production builds embed the compiled frontend into the binary (via go:embed in cmd/testapp/frontend).

### Acknowledgements
- Echo web framework (github.com/labstack/echo)
- go-playground/validator for validation
- fxamacker/cbor for CBOR serialization
- Bun for frontend bundling

### License
 
This project is offered under a dual license:
 
- **Open Source:** You may use, modify, and distribute this software under the terms of the [GNU Affero General Public License v3.0 (AGPL-3.0)](https://www.gnu.org/licenses/agpl-3.0.html).
- **Commercial:** If you would like to use this software under different terms (for example, without the obligations of the AGPL), please [contact us](mailto:info@molpe.de) to obtain a commercial license.

This project includes third-party code components:

1. cbor-js  
   Modified Source: https://github.com/molpeDE/spark/blob/main/internal/typescript/cbor.ts  
   Original Source: https://github.com/paroga/cbor-js  
   License: MIT 

2. typescriptify-golang-structs  
   Modified Source: https://github.com/molpeDE/spark/blob/main/internal/tsgen/tsgen.go  
   Original Source: https://github.com/tkrajina/typescriptify-golang-structs  
   License: Apache-2.0  

All such components retain their original licenses and are NOT covered
by the dual licensing of this project.
