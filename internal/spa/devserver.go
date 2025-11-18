//go:build debug

package spa

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"

	"github.com/molpeDE/spark/internal/typescript"
)

func SPA(_ fs.ReadDirFS, types string) http.Handler {

	var _, callerSource, _, _ = runtime.Caller(2)

	frontendDir := path.Dir(callerSource)

	os.WriteFile(path.Join(frontendDir, "gotypes.ts"), []byte(types), os.ModePerm)

	var ctx, _ = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	if _, err := os.Stat(path.Join(frontendDir, "spark")); err != nil {
		exec.Command("cp", "-R", path.Join(typescript.Dir, "lib"), fmt.Sprintf("%s/spark", frontendDir)).Run()
		exec.Command("ln", "-s", path.Join(typescript.Dir, "bundler.ts"), fmt.Sprintf("%s/spark/bundler", frontendDir)).Run()
	}

	var bundler = exec.CommandContext(ctx, "bun", path.Join(typescript.Dir, "bundler.ts"))
	bundler.Dir = frontendDir
	bundler.Stdout = os.Stdout
	bundler.Stderr = os.Stderr

	bundler.Start()

	var _url, _ = url.Parse("http://localhost:5173")
	return httputil.NewSingleHostReverseProxy(_url)
}
