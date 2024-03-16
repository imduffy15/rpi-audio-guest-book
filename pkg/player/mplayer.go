package player

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type MPlayer struct {
	Volume  int
	lock    sync.Mutex
	process *os.Process
}

func (m *MPlayer) Shutdown() {
	if m.process != nil {
		_ = m.process.Signal(os.Interrupt)
		m.process = nil
	}
}

func (m *MPlayer) Play(fileName string) (err error) {
	if m.lock.TryLock() {
		defer m.lock.Unlock()
		defer func() {
			m.process = nil
		}()

		cmd := exec.Command("mplayer", "-volume", strconv.Itoa(m.Volume), "-cache", "8092", "-", fileName)

		fmt.Printf("Running %s\n", cmd)

		err = cmd.Start()
		m.process = cmd.Process
		if err != nil {
			return
		}

		return cmd.Wait()
	} else {
		return fmt.Errorf("player is busy")
	}
}
