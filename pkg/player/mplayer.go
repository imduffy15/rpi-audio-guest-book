package player

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type MPlayer struct {
	Volume  int
	lock    sync.Mutex
	process *os.Process
}

func (m *MPlayer) Shutdown() {
	fmt.Printf("Player is shutting down\n")
	if m.process != nil {
		_ = m.process.Signal(os.Interrupt)
		m.process = nil
	}
	for !m.lock.TryLock() {
		fmt.Println("Trying to unlock from shutdown function.")
		time.Sleep(1000 * time.Millisecond)
	}
	defer m.lock.Unlock()
}

func (m *MPlayer) Play(fileName string) (err error) {
	if m.lock.TryLock() {
		defer m.lock.Unlock()
		defer func() {
			m.process = nil
		}()

		cmd := exec.Command("mplayer", "-volume", strconv.Itoa(m.Volume), "-cache", "8092", "-", fileName)

		var stderrPipe io.ReadCloser

		// stdoutPipe, err = cmd.StdoutPipe()
		// if err != nil {
		// 	fmt.Println("Error setting up stdout pipe:", err)
		// 	return
		// }

		stderrPipe, err = cmd.StderrPipe()
		if err != nil {
			fmt.Println("Error setting up stderr pipe:", err)
			return
		}

		err = cmd.Start()
		m.process = cmd.Process
		if err != nil {
			fmt.Println("Error starting command:", err)
			return
		}

		var stderrBuf bytes.Buffer
		// go func() {
		// 	io.Copy(&stdoutBuf, stdoutPipe)
		// }()
		go func() {
			io.Copy(&stderrBuf, stderrPipe)
		}()

		err = cmd.Wait()
		if err != nil {
			fmt.Println("Error waiting for command:", err)
		}

		// stdoutStr := stdoutBuf.String()
		stderrStr := stderrBuf.String()

		// fmt.Println("Stdout:", stdoutStr)
		fmt.Println("Stderr:", stderrStr)
		return err
	} else {
		fmt.Printf("Player is busy\n")
		return err
	}

}
