//go:build darwin
// +build darwin

package input_listener

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/imduffy15/rpi-audio-guest-book/pkg/telephone"
)

type GPIOInputListener struct {
	ctx            context.Context
	isShuttingDown atomic.Bool
}

func NewGPIOInputListener(gpioboard string, gpiopin int) *GPIOInputListener {
	return &GPIOInputListener{}
}

func (g *GPIOInputListener) Start(telephone *telephone.Telephone) error {
	var cancel context.CancelFunc
	g.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	return fmt.Errorf("GPIO input listener not supported on darwin")
}

func (g *GPIOInputListener) Shutdown() {
	g.isShuttingDown.Store(true)
	<-g.ctx.Done()
}
