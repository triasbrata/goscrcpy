package path

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/triasbrata/goscrcpy/internals/cons"
)

var adbPath = ""
var serverLocalPath = ""

func GetServerLocalPath() string {
	if len(serverLocalPath) != 0 {
		return serverLocalPath
	}
	serverLocalPath = path.Join(GetExePath(), cons.SERVER_LOCAL_PATH)
	_, err := os.Stat(serverLocalPath)
	if err != nil && !os.IsExist(err) {
		panic(fmt.Errorf("error when get adb path: %w", err))
	}

	return serverLocalPath
}
func GetAdbPath() string {
	if len(adbPath) != 0 {
		return adbPath
	}
	adbPathOs := cons.ADB_PATH_WINDOWS
	if runtime.GOOS == "linux" {
		adbPathOs = cons.ADB_PATH_WINDOWS
	}
	absAdbPath, _ := filepath.Abs(adbPathOs)
	adbPath = absAdbPath
	if _, err := os.Stat(adbPath); err != nil && !os.IsExist(err) {
		panic(fmt.Errorf("error when get adb path: %w", err))
	}
	return adbPath
}
func GetExePath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("error when create client:\n %w", err))
	}
	return dir
}
