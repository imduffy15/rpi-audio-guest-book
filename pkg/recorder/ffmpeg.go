package recorder

import (
	"bytes"
	"fmt"
	"io"
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
	fmt.Printf("Recorder is shutting down\n")
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

		cmd := exec.Command("ffmpeg", "-f", "alsa", "-channels", "1", "-ac", "1", "-i", "hw:3", fmt.Sprintf("%s%s", f.RecordingsPath, filename))
		fmt.Printf("Running %s\n", cmd)

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
		f.process = cmd.Process
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
		fmt.Printf("Recorder is busy\n")
		return fmt.Errorf("recorder is busy")
	}
}
