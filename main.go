package main

import (
	"context"
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	viperConfig "github.com/imduffy15/rpi-audio-guest-book/pkg/config"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/input_listener"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/player"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/recorder"
	"github.com/imduffy15/rpi-audio-guest-book/pkg/telephone"
)

// Version is set during build.
var Version = "dev"

func main() {
	ctx := context.Background()
	config := viperConfig.LoadViperConfig()
	if err := run(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, config *viperConfig.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Required ffmpeg to be installed
	recorder := &recorder.FFmpeg{RecordingsPath: config.RecordingsPath}
	defer recorder.Shutdown()

	// Required mplayer to be installed
	player := &player.MPlayer{Volume: config.Volume}
	defer player.Shutdown()

	// This provider TTS, on first run they are generated/downloaded and cached to ./audio
	speech := htgotts.Speech{Folder: "audio", Language: voices.EnglishUK, Handler: player}

	telephone := telephone.NewTelephone(speech, player, recorder)

	var inputListener input_listener.InputListener
	if config.InputListener == viperConfig.Keyboard {
		inputListener = input_listener.NewKeyboardInputListener()
	} else if config.InputListener == viperConfig.GPIO {
		inputListener = input_listener.NewGPIOInputListener(config.GPIOBoard, config.GPIOPin)
	} else {
		return fmt.Errorf("invalid input listener: %s", config.InputListener)
	}
	defer inputListener.Shutdown()

	inputListenerErr := make(chan error, 1)
	go func() {
		inputListenerErr <- inputListener.Start(telephone)
	}()

	fs := http.FileServer(http.Dir(config.RecordingsPath))
	http.Handle("/", http.StripPrefix("/", fs))
	go func() {
		_ = http.ListenAndServe(":80", nil)
	}()

	select {
	case err := <-inputListenerErr:
		if err != nil {
			return err
		}
		cancel()
	case <-ctx.Done():
		cancel()
	}

	return nil
}
