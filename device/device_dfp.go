//go:build dfp

package device

import (
	"bytes"
	"encoding/binary"
	"log"
	"sync"
	"time"
)

const (
	vid = 0x046d
	pid = 0xc298
)

// ff1f 00 80 80 ffad3f
type state struct {
	Axis     uint16
	Btn      byte
	Hat      byte
	_        byte
	Throttle byte
	Brake    byte
	Dummy    byte
}

func (s *state) State() State {
	return State{
		Axises: []float32{
			float32(int16(s.Axis&0x3fff)-0x2000) / 0x2000,
			float32(s.Throttle^0xff) / 255,
			float32(s.Brake^0xff) / 255,
		},
		Buttons: []bool{
			s.Axis&0x8000 != 0,
			s.Axis&0x4000 != 0,
			s.Btn&0x01 != 0,
			s.Btn&0x02 != 0,
			s.Btn&0x04 != 0,
			s.Btn&0x08 != 0,
			s.Btn&0x10 != 0,
			s.Btn&0x20 != 0,
			s.Btn&0x40 != 0,
			s.Btn&0x80 != 0,
			s.Hat&0x01 != 0,
			s.Hat&0x02 != 0,
			s.Hat&0x04 != 0,
			s.Hat&0x08 != 0,
		},
		Hats: []byte{s.Hat >> 4},
	}
}

var (
	setup = [][]byte{
		{0, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		{0, 0xfe, 0x0d, 0x0c, 0x0c, 0x80, 0x00, 0x00, 0x00},
		{0, 0x11, 0x08, 0x80, 0x80, 0x00, 0x00, 0x00, 0x00},
		{0, 0x21, 0x0c, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00},
		{0, 0x05, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
)

type dev struct {
	base
	mu    sync.RWMutex
	state state
}

func (d *dev) init() error {
	for _, v := range setup {
		time.Sleep(100 * time.Millisecond)
		if _, err := d.device.Write(v); err != nil {
			return err
		}
	}
	return nil
}

func (d *dev) update(b []byte) {
	var s state
	if err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &s); err != nil {
		log.Println(err)
		return
	}
	d.mu.Lock()
	d.state = s
	d.mu.Unlock()
}

func (d *dev) State() State {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.State()
}

func (d *dev) Force(f float32) {
	v := int16(128*f) + 128
	switch {
	case v > 255:
		v = 255
	case v < 0:
		v = 0
	}
	if _, err := d.device.Write([]byte{0, 0x11, 0x08, byte(v), 0x80, 0x00, 0x00, 0x00}); err != nil {
		log.Println(err)
		return
	}
}
