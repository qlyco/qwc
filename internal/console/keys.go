package console

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/qlyco/qwc/core"
	"github.com/qlyco/qwc/internal/res"
)

func Listen(server *core.Server, console Terminal) {
	go stopOnSignal(console)
	go readKeys(func(key byte) {
		handleKey(server, console, key)
	})
}

func handleKey(server *core.Server, console Terminal, key byte) {
	switch key {
	case res.ReloadKey:
		server.Reload()
	case res.RefreshKey:
		server.Refresh()
	case res.InterruptKey:
		stop(console)
	}
}

func stopOnSignal(console Terminal) {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	stop(console)
}

func stop(console Terminal) {
	console.Restore()
	os.Exit(res.StopExitCode)
}

func readKeys(onKey func(key byte)) {
	key := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(key)

		if err != nil {
			return
		}

		onKey(key[0])
	}
}
