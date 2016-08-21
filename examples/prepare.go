package examples

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

func testPrepare() {
	_, filename, _, _ := runtime.Caller(1)
	dir := path.Dir(filename)
	fmt.Println(dir)

	os.Remove(path.Join(dir, "rbac.db"))
}
