//go:build !dfp && !dummy

package device

import (
	"encoding/binary"
	"log"
)

const (
	vid = 0x2e8a
	pid = 0x000a
)

type dev struct {
	base
	state State
}

func (d *dev) init() error {
	if _, err := d.device.Write([]byte{0x0c, 0x04}); err != nil {
		return err
	}
	if _, err := d.device.Write([]byte{0x0c, 0x03}); err != nil {
		return err
	}
	if _, err := d.device.Write([]byte{0x0c, 0x01}); err != nil {
		return err
	}
	if _, err := d.device.Write([]byte{0x00, 0x01, 0x00, 0x00}); err != nil {
		return err
	}
	if _, err := d.device.Write([]byte{
		0x01, 0x01, 0x01, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x04, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00,
	}); err != nil {
		return err
	}
	if _, err := d.device.Write([]byte{0x0a, 0x01, 0x01, 0x01}); err != nil {
		return err
	}
	d.state.Axises = []float32{0}
	return nil
}

func (d *dev) update(b []byte) {
	v := float32(int16(binary.LittleEndian.Uint16(b[0:2]))) / 32767
	switch {
	case v < -1:
		v = -1
	case v > 1:
		v = 1
	}
	d.state.Axises[0] = v
}

func (d *dev) State() State {
	return d.state
}

func (d *dev) Force(f float32) {
	v := int16(32768*f) + 128
	switch {
	case v > 32767:
		v = 32767
	case v < -32767:
		v = -32767
	}
	if _, err := d.device.Write([]byte{0x05, 0x01, byte(v & 0xff), byte(v >> 8)}); err != nil {
		log.Println(err)
		return
	}
	if _, err := d.device.Write([]byte{
		0x01, 0x01, 0x01, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x04, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00,
	}); err != nil {
		log.Println(err)
		return
	}
	if _, err := d.device.Write([]byte{0x0a, 0x01, 0x01, 0x01}); err != nil {
		log.Println(err)
		return
	}
}
