package frontend

import (
	"embed"
	"io/fs"

	"github.com/molpeDE/spark/cmd/testapp/app"
)

//go:generate bun node_modules/appbase/bundler --prod
//go:embed all:dist
var _distEmbed embed.FS
var _dist fs.ReadDirFS

func init() {
	var subfs, _ = fs.Sub(_distEmbed, "dist")
	_dist = subfs.(fs.ReadDirFS)
}

var SPA = app.Instance.SPA(_dist)
