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

	"github.com/warthog618/go-gpiocdev"
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

func (g *GPIOInputListener) Start(telephone *telephone.Telephone) error {
	var cancel context.CancelFunc
	g.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("Input listener (%s) started...\n", viperConfig.GPIO)

	l, err := gpiocdev.RequestLine(g.gpioBoard, g.gpioPin,
		gpiocdev.WithPullUp,
		gpiocdev.WithFallingEdge,
		gpiocdev.WithEventHandler(func(event gpiocdev.LineEvent) {
			if event.Type == gpiocdev.LineEventFallingEdge {
				telephone.ToggleState()
			}
		}))
	if err != nil {
		return fmt.Errorf("failed to start gpio input listener: %w", err)
	}
	defer l.Close()

	for {
		if g.isShuttingDown.Load() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (g *GPIOInputListener) Shutdown() {
	g.isShuttingDown.Store(true)
	<-g.ctx.Done()
}
