package engine

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

func capture(w io.Writer, r io.Reader) ([]byte, error) {
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

// ExecuteAndCapture function
func ExecuteAndCapture(path string, args []string) (<-chan *exec.Cmd, error) {
	out := make(chan *exec.Cmd)
	cmd := exec.Command(path, args...)
	stdoutIn, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = capture(os.Stdout, stdoutIn)
		wg.Done()
	}()
	wg.Wait()
	out <- cmd
	return out, nil
	//err = cmd.Wait()
	//if err != nil {
	//	return cmd, err
	//}
}

// KillCommand function
func KillCommand(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}
