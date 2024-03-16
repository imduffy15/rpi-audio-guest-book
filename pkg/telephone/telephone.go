package telephone

import (
	"fmt"
	"sync/atomic"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/player"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/recorder"
)

type State int

const (
	OnHook State = iota
	OffHook
)

type Telephone struct {
	speech   htgotts.Speech
	player   player.Player
	recorder recorder.Recorder
	state    atomic.Value
}

func NewTelephone(speech htgotts.Speech, player player.Player, recorder recorder.Recorder) *Telephone {
	telephone := &Telephone{
		speech:   speech,
		player:   player,
		recorder: recorder,
	}
	telephone.state.Store(OnHook)
	return telephone
}

func (t *Telephone) ToggleState() {
	currentState := t.state.Load().(State)
	if currentState == OnHook {
		t.transition(OffHook)
	} else {
		t.transition(OnHook)
	}
}

type transition func(t *Telephone)

var transitionTable = map[State]transition{
	OnHook: func(t *Telephone) {
		fmt.Printf("Phone is on the hook\n")
		t.player.Shutdown()
		t.recorder.Shutdown()
	},
	OffHook: func(t *Telephone) {
		fmt.Printf("Phone is off the hook\n")

		err := t.speech.Speak("Please leave your greeting after the beep, hang up when your finished.")
		if err != nil {
			return
		}
		defer t.player.Shutdown()

		recording := make(chan error, 1)
		go func() {

			recording <- t.recorder.Record(fmt.Sprintf("%d.mp3", time.Now().Unix()))
		}()
		defer t.recorder.Shutdown()

		err = t.player.Play("audio/beep.mp3")
		if err != nil {
			return
		}
		defer t.player.Shutdown()

		select {
		case <-recording:
		case <-time.After(1 * time.Minute):
			t.recorder.Shutdown()
		}
	},
}

func (t *Telephone) transition(newState State) {
	if transitionFunc, ok := transitionTable[newState]; ok {
		go transitionFunc(t)
		t.state.Store(newState)
	}
}
