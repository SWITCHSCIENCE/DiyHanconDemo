//go:build dummy

package device

import (
	"context"
	"log"
	"math"
	"time"
)

const (
	MaxToruque = 9.7                          // Nm
	Inertia    = 0.01                         // kg.m^2 直径20cm、厚み5cm、重さ2Kgの円柱のイナーシャ
	L2Ldeg     = float32(540 * math.Pi / 180) // 540deg to rad
)

type dev struct {
	state    State
	force    float32
	velocity float32 // rad/s
	position float32 // rad
	tm       time.Time
}

func (d *dev) init() error {
	d.state.Axises = []float32{0}
	return nil
}

func (d *dev) update() {
	defer func() { d.tm = time.Now() }()
	if d.tm.IsZero() {
		return
	}
	elapsed := float32(time.Since(d.tm).Seconds())
	acceleration := MaxToruque * d.force / Inertia
	d.velocity = d.velocity + acceleration*elapsed - d.position/100
	d.position = d.position + d.velocity*elapsed
	switch {
	case d.position > L2Ldeg:
		d.position = L2Ldeg
	case d.position < -L2Ldeg:
		d.position = -L2Ldeg
	}
	d.state.Axises[0] = d.position / L2Ldeg // rad to v
}

func (d *dev) Close() error {
	return nil
}

func (d *dev) State() State {
	return d.state
}

func (d *dev) Force(f float32) {
	d.force = f
}

func New(ctx context.Context) Device {
	d := &dev{}
	if err := d.init(); err != nil {
		log.Fatal(err)
	}
	go func() {
		defer d.Close()
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				d.update()
			}
		}
	}()
	return d
}
