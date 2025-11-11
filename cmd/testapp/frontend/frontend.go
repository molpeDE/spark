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
