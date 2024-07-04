package telephone

import (
	"fmt"
	"sync"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/player"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/recorder"
	"time"
	"sync/atomic"
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
	state atomic.Value
	wg       *sync.WaitGroup
}

func NewTelephone(speech htgotts.Speech, player player.Player, recorder recorder.Recorder) *Telephone {
	telephone := &Telephone{
		speech:   speech,
		player:   player,
		recorder: recorder,
		wg:       &sync.WaitGroup{},
	}
	return telephone
}

type transition func(t *Telephone) error

var transitionTable = map[State]transition{
	OnHook: func(t *Telephone) error {
		defer t.wg.Done()
		fmt.Printf("Phone is on the hook\n")
		t.player.Shutdown()
		t.recorder.Shutdown()
		return nil
	},
	OffHook: func(t *Telephone) error {
		fmt.Printf("Phone is off the hook\n")

		fmt.Printf("Playing the greeting\n")
		err := t.speech.Speak("Please leave your greeting after the beep, hang up when your finished.")
		if err != nil {
			fmt.Println("Greeting playback errored", err)
			return err
		}

		fmt.Printf("Recording...\n")
		recording := make(chan error, 1)
		go func() {
			recording <- t.recorder.Record(fmt.Sprintf("%d.wav", time.Now().Unix()))
		}()

		fmt.Printf("Playing the beep\n")
		err = t.player.Play("audio/beep.mp3")
		if err != nil {
			fmt.Println("Beep playback errored", err)
			return err
		}

		fmt.Printf("Waiting for OffHook to finish...\n")
		select {
		case <-recording:
		case <-time.After(1 * time.Minute):
			fmt.Printf("Timed out, stopping\n")
		}

		return nil
	},
}

func (t *Telephone) Transition(newState State) {
	if newState == t.state.Load() && newState == OffHook {
		t.wg.Add(1)
		transitionTable[OnHook](t)
	}

	if transitionFunc, ok := transitionTable[newState]; ok {
		if newState == OffHook {
			fmt.Printf("Waiting for OnHook operations to finish\n")
			t.wg.Wait()
		} else {
			t.wg.Add(1)
		}
		fmt.Printf("Transitioning to %s\n", newState)
		go transitionFunc(t)
		t.state.Store(newState)
	}
}
