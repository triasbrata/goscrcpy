package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/triasbrata/gadb"
)

type AdbCmd struct {
	adbPath       string
	port          int
	ipHost        string
	client        *gadb.Client
	serverRunning bool
}

func (adb *AdbCmd) run(args ...string) (string, error) {
	cmd := exec.Command(adb.adbPath, args...)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	res, err := cmd.Output()
	return string(res), err
}
func (adb *AdbCmd) NewServer() error {
	_, err := adb.run("-L", fmt.Sprintf("tcp:%v", adb.port), "start-server")
	if err != nil && !strings.Contains(err.Error(), "already started") {

		return err
	}
	adb.serverRunning = true
	client, err := gadb.NewClientWith(adb.ipHost, adb.port)
	if err != nil {
		return err
	}
	adb.client = &client
	return nil
}
func (adb *AdbCmd) KillServer() error {
	if !adb.serverRunning {
		return fmt.Errorf("server not running yet")
	}
	return adb.client.KillServer()

}
func (adb *AdbCmd) ListDevices() ([]gadb.Device, error) {
	if !adb.serverRunning {
		return nil, fmt.Errorf("server not running yet")
	}
	return adb.client.DeviceList()
}

func NewAdbExe(adbPath string) *AdbCmd {
	return &AdbCmd{
		adbPath: adbPath,
		port:    5037,
		ipHost:  "localhost",
	}
}
