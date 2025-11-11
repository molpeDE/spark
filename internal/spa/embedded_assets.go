//go:build !debug

package spa

import (
	"io/fs"
	"net/http"
	"path"
)

func SPA(dist fs.ReadDirFS, _ string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check whether a file exists or is a directory at the given path
		// if file does not exist or path is a directory, serve index.html
		if fi, err := fs.Stat(dist, path.Join(".", r.URL.Path)); err != nil || fi.IsDir() {
			http.ServeFileFS(w, r, dist, "index.html")
			return
		}

		// otherwise, serve the static file
		http.FileServerFS(dist).ServeHTTP(w, r)
	})
}
