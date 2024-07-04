package input_listener

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
	viperConfig "github.com/imduffy15/rpi-audio-guest-book/pkg/config"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/telephone"
)

type KeyboardInputListener struct {
	ctx            context.Context
	isShuttingDown atomic.Bool
}

func NewKeyboardInputListener() *KeyboardInputListener {
	return &KeyboardInputListener{}
}

func (k *KeyboardInputListener) Start(t *telephone.Telephone) error {
	var cancel context.CancelFunc
	k.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		return err
	}
	defer keyboard.Close()

	fmt.Printf("Input listener (%s) started...\n", viperConfig.Keyboard)
	for {
		if k.isShuttingDown.Load() {
			return nil
		}

		event := <-keysEvents
		if event.Err != nil {
			return event.Err
		}
		if event.Key == keyboard.KeySpace {
			t.Transition(telephone.OffHook)
		} else if event.Key == keyboard.KeyCtrlC {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (k *KeyboardInputListener) Shutdown() {
	k.isShuttingDown.Store(true)
	<-k.ctx.Done()
}
