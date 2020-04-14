package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-qamel/qamel"
)

type Backend struct {
	qamel.QmlObject
	_ func()                        `constructor:"init"`
	_ func(bool)                    `signal:"changeSearch"`
	_ func(int, int, int, int, int) `signal:"changeDPI"`
	_ func(int, int, int, int, int) `slot:"setDPI"`
	_ func(string, int, int)        `slot:"setColor"`
}

var backend *Backend

func (b *Backend) init() {
	backend = b
}

func (b *Backend) setDPI(dpi1, dpi2, dpi3, dpi4, dpi5 int) {
	d := [][]byte{decode(dpi1, 2), decode(dpi2, 2), decode(dpi3, 2), decode(dpi4, 2), decode(dpi5, 2)}
	dpiData := []byte{}
	for _, data := range d {
		dpiData = append(dpiData, data[1], data[0])
	}

	connect()

	controlTransfer(PreSetDPI.Part1)
	_, _, err := device.BulkTransferIn(0x82, 20, 1000)
	if err != nil {
		log.Println(err)
	}
	controlTransfer(
		append(
			append(SetDPI.Part1, dpiData...),
			SetDPI.Part2...,
		),
	)
	_, _, err = device.BulkTransferIn(0x82, 20, 1000)
	if err != nil {
		log.Println(err)
	}
	controlTransfer(PostSetDPI.Part1)
	_, _, err = device.BulkTransferIn(0x82, 20, 1000)
	if err != nil {
		log.Println(err)
	}

	disconnect()
}

func (b *Backend) setColor(color string, speed, effect int) {
	c := strings.TrimPrefix(color, "#")
	col, _ := strconv.ParseInt(c, 16, 32)
	speed = 13000 - speed

	switch effect {
	case 0:
		controlTransfer(
			append(
				append(SetColor.Part1, decode(int(col), 3)...),
				SetColor.Part2...,
			),
		)
	case 1:
		controlTransfer(
			append(
				append(
					append(SetBreathe.Part1, decode(int(col), 3)...),
					decode(speed, 2)...,
				),
				SetBreathe.Part2...,
			),
		)
	case 2:
		controlTransfer(
			append(
				append(SetCycle.Part1, decode(speed, 2)...),
				SetCycle.Part2...,
			),
		)
	}
}
