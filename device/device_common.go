package device

import (
	"fmt"
	"strings"

	"github.com/sstallion/go-hid"
)

type State struct {
	Buttons []bool
	Hats    []byte
	Axises  []float32
}

func (s State) String() string {
	bits := []string{}
	for _, v := range s.Buttons {
		if v {
			bits = append(bits, "1")
		} else {
			bits = append(bits, "0")
		}
	}
	return fmt.Sprintf("axises:%v buttons:%v hats:%v", s.Axises, strings.Join(bits, ""), s.Hats)
}

type Device interface {
	Close() error
	State() State
	Force(f float32)
}

type base struct {
	device *hid.Device
}

func (d *base) Close() error {
	return d.device.Close()
}
