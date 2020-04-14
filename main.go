package main

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-qamel/qamel"
	"github.com/gotmc/libusb"
)

type Command struct {
	Part1 []byte
	Part2 []byte
}

// hot reload window and debug logs
const debug = false

var (
	ctx    *libusb.Context
	device *libusb.DeviceHandle

	connected bool

	vendor  uint16 = 0x046d
	product uint16 = 0xc084
	wIndex         = 0x0001

	req byte   = 0x09
	val uint16 = 0x0210

	RequestDPI = Command{
		[]byte{0x11, 0xff, 0x0f, 0x5a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{},
	}
	PreSetDPI = Command{
		[]byte{0x11, 0xff, 0x0f, 0x6d, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{},
	}
	PostSetDPI = Command{
		[]byte{0x10, 0xff, 0x0f, 0x8d, 0x00, 0x00, 0x00}, []byte{},
	}
	SetDPI = Command{
		[]byte{0x11, 0xff, 0x0f, 0x7d, 0x01, 0x02, 0x02}, /* 5x2 dpi value in bytes (eg 2C 01 E8 03 6C 07 94 11 00 00), */
		[]byte{0xff, 0xff, 0xff},                         /* !WARNING! dpi is "reversed" so 300 = 0x012c in hex = []byte{0x2c, 0x01} as dpi */
	}
	SetColor = Command{
		[]byte{0x11, 0xff, 0x0e, 0x3c, 0x00, 0x01}, /* 3 color bytes (eg 0xff, 0x00, 0x00), */
		[]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	SetBreathe = Command{
		[]byte{0x11, 0xff, 0x0e, 0x3c, 0x00, 0x03}, /* 3 color bytes (eg 0xff, 0x00, 0x00), */
		/* 2 speed bytes (eg 0x13, 0x88), */ []byte{0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	SetCycle = Command{
		[]byte{0x11, 0xff, 0x0e, 0x3c, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00},
		/* 2 speed bytes (eg 0x13, 0x88), */ []byte{0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
)

func init() {
	// Register the Backend as QML component
	RegisterQmlBackend("Backend", 1, 0, "Backend")
}

func main() {
	// Create application
	app := qamel.NewApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("G203 Control")

	// Create a QML viewer
	var view qamel.Viewer
	if debug {
		view = qamel.NewViewerWithSource("res/main.qml")

		projectDir, err := os.Getwd()
		if err != nil {
			log.Fatalln("Failed to get working directory:", err)
		}

		resDir := filepath.Join(projectDir, "res")
		go view.WatchResourceDir(resDir)
	} else {
		view = qamel.NewViewerWithSource("qrc:/res/main.qml")
	}
	view.SetResizeMode(qamel.SizeRootObjectToView)
	view.SetMinimumWidth(400)
	view.SetMaximumWidth(400)
	view.SetMinimumHeight(300)
	view.SetMaximumHeight(300)
	view.Show()

	go func() {
		connect()

		// Request DPI Data from mouse
		controlTransfer(RequestDPI.Part1)
		dpiData, _, err := device.BulkTransferIn(0x82, 20, 1000)
		if err != nil {
			log.Fatalf("Failed to read DPI Data: %v\n", err)
		}
		log.Printf("DPI Data: %v\n", dpiData)
		backend.changeDPI(
			encode([]byte{dpiData[8], dpiData[7]}),
			encode([]byte{dpiData[10], dpiData[9]}),
			encode([]byte{dpiData[12], dpiData[11]}),
			encode([]byte{dpiData[14], dpiData[13]}),
			encode([]byte{dpiData[16], dpiData[15]}),
		)
		backend.changeSearch(false)

		disconnect()
	}()

	// Exec app
	app.Exec()
}

func connect() {
	ctx, _ = libusb.NewContext()

	var err error
	_, device, err = ctx.OpenDeviceWithVendorProduct(vendor, product)
	if err != nil {
		log.Fatalf("Could not open a device: %v\nTry to run as root (linux; with sudo) or run as administrator (windows; right click)", err)
	}

	f, _ := device.KernelDriverActive(wIndex)
	if !f {
		device.AttachKernelDriver(wIndex)
	}
	device.DetachKernelDriver(wIndex)
	log.Println("Connected")
	connected = true
}

func disconnect() {
	err := device.AttachKernelDriver(wIndex)
	if err != nil {
		log.Printf("[warn] Failed to attach kernel driver: %v\n", err)
	}
	err = device.Close()
	if err != nil {
		log.Printf("Failed to close device handle: %v\n", err)
	}
	err = ctx.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Disconnected")
	connected = false
}

func controlTransfer(data []byte) {
	if debug {
		log.Printf("Control transfer: %v\n", data)
	}
	f := false
	if !connected {
		connect()
		f = true
	}
	_, err := device.ControlTransfer(0x21, req, val, uint16(wIndex), data, 20, 10)
	if err != nil {
		log.Println(err)
	}
	if f {
		disconnect()
	}
}

func encode(data []byte) int {
	i64, err := strconv.ParseInt(hex.EncodeToString(data), 16, 64)
	if err != nil {
		log.Printf("Failed to encode data: %v\n", err)
	}
	return int(i64)
}

func decode(data, l int) (d []byte) {
	s := strconv.FormatInt(int64(data), 16)
	if len(s)%2 != 0 {
		s = "0" + s
	}
	d, err := hex.DecodeString(s)
	if err != nil {
		log.Printf("Failed to decode data: %v\n", err)
	}
	for len(d) < l {
		d = append([]byte{0}, d...)
	}
	return
}
