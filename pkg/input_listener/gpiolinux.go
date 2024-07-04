//go:build linux
// +build linux

package input_listener

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	viperConfig "github.com/imduffy15/rpi-audio-guest-book/pkg/config"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/telephone"

	"github.com/stianeikeland/go-rpio/v4"
)

type GPIOInputListener struct {
	ctx            context.Context
	isShuttingDown atomic.Bool
	gpioBoard      string
	gpioPin        int
}

func NewGPIOInputListener(gpioboard string, gpiopin int) *GPIOInputListener {
	return &GPIOInputListener{gpioPin: gpiopin}
}

func (g *GPIOInputListener) Start(t *telephone.Telephone) error {
	var cancel context.CancelFunc
	g.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("Input listener (%s) started...\n", viperConfig.GPIO)

	err := rpio.Open()
	if err != nil {
		return fmt.Errorf("failed to start gpio input listener: %w", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(g.gpioPin)
	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge)
	pin.Detect(rpio.RiseEdge)
	defer pin.Detect(rpio.NoEdge)

	previousState := pin.Read()

	go func() {
		for {
			if pin.EdgeDetected() {
				currentState := pin.Read()
				fmt.Printf("Edge detcted, previous state: %s current state: %s\n", previousState, currentState)
				if currentState == rpio.Low {
					t.Transition(telephone.OffHook)
				} else {
					t.Transition(telephone.OnHook)
				}
				previousState = currentState

			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for {
		if g.isShuttingDown.Load() {
			return nil
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (g *GPIOInputListener) Shutdown() {
	g.isShuttingDown.Store(true)
	<-g.ctx.Done()
}
