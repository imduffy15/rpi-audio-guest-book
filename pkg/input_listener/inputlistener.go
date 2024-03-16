package input_listener

import "github.com/imduffy15/rpi-audio-guest-book/pkg/telephone"

type InputListener interface {
	Start(t *telephone.Telephone) error
	Shutdown()
}
