package reader

import (
	"fmt"
	"io"
	"time"
)

type DelayedReadWriter struct {
	startDelay  time.Duration
	middleDelay time.Duration
	endDelay    time.Duration

	reader io.ReadCloser
	writer io.WriteCloser
}

func NewDelayedReadWriter(s, m, e int) *DelayedReadWriter {
	r, w := io.Pipe()

	return &DelayedReadWriter{
		startDelay:  time.Duration(s),
		middleDelay: time.Duration(m),
		endDelay:    time.Duration(e),

		reader: r,
		writer: w,
	}
}

func (drw *DelayedReadWriter) Read(p []byte) (int, error) {
	return drw.reader.Read(p)
}

func (drw *DelayedReadWriter) Write(msg string) error {
	time.Sleep(drw.startDelay * time.Second)
	_, err := drw.writer.Write([]byte(msg)[0 : len(msg)/2])
	if err != nil {
		return fmt.Errorf("Error writing part 1: %s\n", err)
	}

	time.Sleep(drw.middleDelay * time.Second)
	_, err = drw.writer.Write([]byte(msg)[len(msg)/2 : len(msg)])
	if err != nil {
		return fmt.Errorf("Error writing part 1: %s\n", err)
	}

	time.Sleep(drw.endDelay * time.Second)

	drw.writer.Close()

	return nil
}
