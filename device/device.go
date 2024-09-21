//go:build !dummy

package device

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/sstallion/go-hid"
)

func New(ctx context.Context) Device {
	res := make(chan *hid.Device, 1)
	hid.Enumerate(vid, pid, func(info *hid.DeviceInfo) error {
		defer close(res)
		log.Printf("%#v", info)
		device, err := hid.Open(info.VendorID, info.ProductID, info.SerialNbr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%#v", device)
		b := make([]byte, 8192)
		n, err := device.GetReportDescriptor(b)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("desc: %X", b[:n])
		res <- device
		return nil
	})
	device := <-res
	if device == nil {
		log.Fatal("no device found")
	}
	d := &dev{
		base: base{device: device},
	}
	if err := d.init(); err != nil {
		log.Fatal(err)
	}
	go func() {
		defer d.Close()
		b := make([]byte, 128)
		for {
			select {
			default:
			case <-ctx.Done():
				return
			}
			n, err := d.device.ReadWithTimeout(b, 20*time.Millisecond)
			if err != nil {
				if errors.Is(err, hid.ErrTimeout) {
					continue
				}
				log.Fatal(err)
			}
			d.update(b[:n])
		}
	}()
	return d
}
