package recorder

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type FFmpeg struct {
	RecordingsPath string
	lock           sync.Mutex
	process        *os.Process
}

func (f *FFmpeg) Shutdown() {
	if f.process != nil {
		_ = f.process.Signal(os.Interrupt)
		f.process = nil
	}

}

func (f *FFmpeg) Record(filename string) (err error) {
	if f.lock.TryLock() {
		defer f.lock.Unlock()
		defer func() {
			f.process = nil
		}()

		// ffmpeg -f alsa -channels 1 -i hw:0,0 -ac 1 
		cmd := exec.Command("ffmpeg", "-f", "alsa", "-channels", "1", "-ac", "1", "-i", "hw:3", fmt.Sprintf("%s%s", f.RecordingsPath, filename))

		fmt.Printf("Running %s\n", cmd)

		err = cmd.Start()
		f.process = cmd.Process
		if err != nil {
			return
		}

		err = cmd.Wait()
		return err
	} else {
		return fmt.Errorf("recorder is busy")
	}
}
