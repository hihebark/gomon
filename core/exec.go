package core

import (
	"errors"
	//"fmt"
	"io"
	//"log"
	"os"
	"os/exec"
	"sync"
)

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func Execute(path string, args []string) (string, error) {
	cmd := exec.Command(path, args...)
	var stdout, stderr []byte
	var errStdout error
	stdoutIn, _ := cmd.StdoutPipe()
	//stderrIn, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	//stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	//if errStdout != nil || errStderr != nil {
	//	return "", errors.New("Empty stdout")
	//}
	return string(stdout), errors.New(string(stderr))
}
