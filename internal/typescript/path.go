package typescript

import (
	"path"
	"runtime"
)

var _, callerSource, _, _ = runtime.Caller(0)
var Dir = path.Dir(callerSource)
