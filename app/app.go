package goscrcpy

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/triasbrata/gadb"
	"github.com/triasbrata/goscrcpy/internals/cmd"
	"github.com/triasbrata/goscrcpy/internals/cons"
	"github.com/triasbrata/goscrcpy/internals/gui"
	"github.com/triasbrata/goscrcpy/internals/path"
	"github.com/triasbrata/goscrcpy/internals/slices"
)

type size struct {
	width  uint16
	height uint16
}
type App struct {
	listener  net.Listener
	deviceCon net.Conn
	adbcmd    *cmd.AdbCmd
}

const DEVICE_NAME_FIELD_LENGTH = 64

func Run() (err error) {
	gadb.SetDebug(true)
	app := App{
		adbcmd: cmd.NewAdbExe(path.GetAdbPath()),
	}

	err = app.adbcmd.NewServer()
	if err != nil {
		return fmt.Errorf("error when create client:\n %w", err)
	}

	defer func() {
		errKill := app.adbcmd.KillServer()
		if errKill != nil {
			fmt.Printf("got error when try kill server %v \n", errKill.Error())
		}
	}()
	// c, _ := adb.New()
	// c.ListDevices()
	devices, err := app.adbcmd.ListDevices()
	if err != nil {
		return err
	}
	fmt.Printf("devices: %v\n", devices)

	devicePaths := slices.Entries(devices, func(it gadb.Device) (string, gadb.Device) {
		return it.Serial(), it
	})
	fmt.Printf("devicePaths: %v\n", devicePaths)
	fmt.Printf("devicePaths: %v\n", devicePaths)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebiten UI - List")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	mainWindow := gui.NewWindow()
	windowSetting := mainWindow.(*gui.WindowSettings)

	//todo need create close aipp and trigger this
	app.bindConfigWindowWithApp(windowSetting)
	if err := ebiten.RunGame(mainWindow); err != nil {
		return err
	}
	return nil

}

func (app *App) connectDevice(device *gadb.Device) {
	var (
		err error
		res string
	)

	// // device := app.adbcmd.Device(adb.DeviceWithSerial(dInfo.Serial))

	appNameServer := "com.genymobile.scrcpy.Server"
	serverVersion := "2.4"
	// localPort := 27183 // reverse时本地监听端口
	// maxSize := 720     // 视频分辨率
	bitRate := 2000000 // 视频比特率
	// maxFps := 0

	//push server
	localServer, err := os.Open(cons.SERVER_LOCAL_PATH)
	if err != nil {
		log.Fatalf("error %w", err)
	}
	err = device.PushFile(localServer, cons.SERVER_REMOTE_PATH)

	if err != nil {
		log.Fatalf("err push command: %v\n", err)
		return
	}

	// start server listener local
	app.listener, err = net.Listen("tcp", ":27183")
	go func() {
		for {
			con, err := app.listener.Accept()
			if err != nil {
				log.Fatalf("error when Accept listener %v", err)
				return
			}
			dn, w, h, err := readInfo(con)
			fmt.Printf("dn: %v\n", dn)
			fmt.Printf("w: %v\n", w)
			fmt.Printf("h: %v\n", h)
			if err != nil {
				log.Fatalf("error when readinfo %v")

			}
		}
	}()
	if err != nil {
		fmt.Printf("err start server: %v\n", err)
		return
	}
	defer func() {
		if app.listener != nil {
			app.listener.Close()
		}
	}()
	err = device.Forward(27183, 27183, false)
	if err != nil {
		log.Fatalf("error when forward port %v", err.Error())
		return
	}
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer func() { wg.Done() }()
	// gadb.DefaultAdbReadTimeout = 30
	//start scrcpy-server
	log.Printf("remote server start")
	res, err = device.RunShellCommandAndForget(fmt.Sprintf("CLASSPATH=%v app_process / %v %v video_bit_rate=%v log_level=verbose lock_video_orientation=1 tunnel_forward=true stay_awake=true audio=false scid=00000f60", cons.SERVER_REMOTE_PATH, appNameServer, serverVersion, bitRate))
	if err != nil {
		log.Fatalf("err run server: %v\n", err)
		return
	}
	if strings.Contains(res, "Aborted") {
		log.Fatalf("abort server")
		return
	}

	log.Println("connect")

	// }()
	// wg.Add(1)
	// go func() {
	// 	defer func() { wg.Done() }()
	// attemps := 100
	// var con net.Conn
	// con, err = app.listener.Accept()
	// if err != nil {
	// 	log.Printf("error when dial con %v \n con: %v", err, con)

	// }
	// for attemps > 0 {
	// 	fmt.Printf("attemps: %v\n", attemps)
	// 	if err != nil {
	// 		log.Printf("error when dial con %v \n con: %v", err, con)
	// 		return
	// 	}
	// 	buf := make([]byte, 1)
	// 	_, err = io.ReadFull(con, buf)
	// 	if err != nil && err.Error() != "EOF" {
	// 		log.Fatalf("error when read buff %v \n con: %v", err, con)
	// 	} else if err != nil {
	// 		err = nil
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// 	attemps--
	// }
	// dn, w, h, err := readInfo(con)
	// fmt.Printf("dn: %v\n", dn)
	// fmt.Printf("sz: %v, %v\n", w, h)
	// fmt.Printf("err: %v\n", err)
	// 	wg.Done()

	// }()
	// wg.Wait()
	// //connect
	// app.deviceCon, err = app.listener.Accept()
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }
	// fmt.Println("read device info")
	// //read device info
	// buf := make([]byte, 68)
	// if _, err = io.ReadFull(app.deviceCon, buf); err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }
	// fmt.Printf("buf: %s\n", buf)

}

func (a *App) bindConfigWindowWithApp(configWindow *gui.WindowSettings) {
	configWindow.SetConnectDevice(a.connectDevice)
	getDeviceInfo := func() {
		serials, err := a.adbcmd.ListDevices()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		configWindow.UpdateSerialList(slices.Map(serials, func(d gadb.Device) *gadb.Device {
			return &d
		}))
	}
	getDeviceInfo()
	configWindow.Cron(getDeviceInfo)
}
func readInfo(videoSocket net.Conn) (string, int, int, error) {
	var buf [DEVICE_NAME_FIELD_LENGTH + 12]byte
	start := time.Now()

	for {
		fmt.Printf("readInfo\n")
		if videoSocket == nil {
			return "", 0, 0, fmt.Errorf("videoSocket is nil")
		}

		videoSocket.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		n, err := videoSocket.Read(buf[:])
		if err, ok := err.(net.Error); ok && err.Timeout() {
			if time.Since(start) > 3*time.Second {
				return "", 0, 0, fmt.Errorf("readInfo timeout")
			}
			continue
		} else if err != nil {
			return "", 0, 0, err
		}

		if n < DEVICE_NAME_FIELD_LENGTH+12 {
			return "", 0, 0, fmt.Errorf("could not retrieve device information")
		}

		break
	}

	buf[DEVICE_NAME_FIELD_LENGTH-1] = 0 // in case the client sends garbage
	deviceName := string(buf[:DEVICE_NAME_FIELD_LENGTH])

	width := int(binary.BigEndian.Uint32(buf[DEVICE_NAME_FIELD_LENGTH+4:]))
	height := int(binary.BigEndian.Uint32(buf[DEVICE_NAME_FIELD_LENGTH+8:]))

	return deviceName, width, height, nil
}
