package demo

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"DiyHanconDemo/device"
)

type PIDController struct {
	Kp, Ki, Kd float32
	prevError  float32
	integral   float32
}

func NewPIDController(kp, ki, kd float32) *PIDController {
	return &PIDController{
		Kp: kp,
		Ki: ki,
		Kd: kd,
	}
}

func (pid *PIDController) Compute(setpoint, process float32, dt float32) float32 {
	e := setpoint - process

	// 比例項
	proportional := pid.Kp * e

	// 積分項
	pid.integral += e * dt
	integral := pid.Ki * pid.integral

	// 微分項
	derivative := pid.Kd * (e - pid.prevError) / dt

	output := proportional + integral + derivative

	pid.prevError = e

	return output
}

type Demo struct {
	device   device.Device
	vibFreq  float64
	vibPower float64
	pid      *PIDController
}

func New(device device.Device) *Demo {
	return &Demo{
		device:   device,
		vibFreq:  5,
		vibPower: 0.2,
		pid:      NewPIDController(0.01, 0.0, 0.0),
	}
}

func (d *Demo) PID(kp, ki, kd float64) {
	d.pid.Kd = float32(kd)
	d.pid.Ki = float32(ki)
	d.pid.Kp = float32(kp)
}

func (d *Demo) Viberation(f, p float64) {
	d.vibFreq = f
	d.vibPower = p
}

func (d *Demo) Run(ctx context.Context) {
	log.Println("demo start")
	defer log.Println("demo end")
	ticker := time.NewTicker(10 * time.Millisecond)
	begin := time.Now()
	last := time.Now()
	for {
		now := time.Now()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dt := float32(now.Sub(last).Seconds())
			state := d.device.State()
			steer := state.Axises[0]
			runtime.EventsEmit(ctx, "steer", steer)
			seconds := now.Sub(begin).Seconds()
			force := float32(d.vibPower * math.Sin(2*math.Pi*d.vibFreq*seconds))
			force += d.pid.Compute(0, steer, dt)
			switch {
			case force > 1:
				force = 1
			case force < -1:
				force = -1
			}
			d.device.Force(force)
			runtime.EventsEmit(ctx, "force", force)
		}
		last = now
	}
}
