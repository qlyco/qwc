package console

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/qlyco/qwc/internal/res"
	"golang.org/x/term"
)

type Terminal struct {
	state *term.State
}

type rawWriter struct {
	destination io.Writer
}

func OpenTerminal() Terminal {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		return Terminal{}
	}

	log.SetOutput(rawWriter{destination: log.Writer()})

	return Terminal{state: state}
}

func (console Terminal) Restore() {
	if console.state != nil {
		term.Restore(int(os.Stdin.Fd()), console.state)
	}
}

func (output rawWriter) Write(text []byte) (int, error) {
	_, err := output.destination.Write(bytes.ReplaceAll(text, []byte(res.LineBreak), []byte(res.RawLineBreak)))

	return len(text), err
}
