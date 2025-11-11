//go:build debug

package spa

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strings"

	"github.com/molpeDE/spark/internal/typescript"
)

func SPA(_ fs.ReadDirFS, types string) http.Handler {

	var _, callerSource, _, _ = runtime.Caller(2)

	frontendDir := path.Dir(callerSource)

	os.WriteFile(path.Join(frontendDir, "gotypes.ts"), []byte(types), os.ModePerm)

	var ctx, _ = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	injectDependency := func() {
		cmd := exec.CommandContext(ctx, "bun", "i", fmt.Sprintf("file:%s", typescript.Dir))
		cmd.Stdout = os.Stdout
		cmd.Dir = frontendDir
		cmd.Stderr = os.Stderr

		// Run the command
		err := cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := os.Stat(path.Join(frontendDir, "package.json")); err != nil {
		injectDependency()
	} else if data, err := os.ReadFile(path.Join(frontendDir, "package.json")); err == nil {
		if !strings.Contains(string(data), typescript.Dir) {
			injectDependency()
		}
	}

	var bundler = exec.CommandContext(ctx, "bun", "node_modules/@molpe/spark/bundler")
	bundler.Dir = frontendDir
	bundler.Stdout = os.Stdout
	bundler.Stderr = os.Stderr
	bundler.Env = append(bundler.Env, fmt.Sprintf("NODE_PATH=%s", path.Join(frontendDir, "node_modules")))

	bundler.Start()

	var _url, _ = url.Parse("http://localhost:5173")
	return httputil.NewSingleHostReverseProxy(_url)
}
