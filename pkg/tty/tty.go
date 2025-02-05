// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source
// code is governed by the MIT license that can be found in the LICENSE
// file.
//
//go:build !windows

package tty

import (
	"os"
	"syscall"

	"github.com/sruehl/term/termios"
)

//======================================================================

type TerminalSignals struct {
	tiosp *syscall.Termios
	out   *os.File
	set   bool
}

func (t *TerminalSignals) IsSet() bool {
	return t.set
}

func (t *TerminalSignals) Restore() {
	if t.out != nil {
		fd := t.out.Fd()
		_ = termios.Tcsetattr(fd, termios.TCSANOW, t.tiosp)

		_ = t.out.Close()
		t.out = nil
	}
	t.set = false
}

func (t *TerminalSignals) Set(outtty string) error {
	var err error
	var newtios syscall.Termios
	var fd uintptr

	if t.out, err = os.OpenFile(outtty, os.O_WRONLY, 0); err != nil {
		goto failed
	}

	fd = t.out.Fd()

	if err = termios.Tcgetattr(fd, t.tiosp); err != nil {
		goto failed
	}

	newtios = *t.tiosp
	newtios.Lflag |= syscall.ISIG

	// Enable ctrl-z for suspending the foreground process group via the
	// line discipline. Ctrl-c and Ctrl-\ are not handled, so the terminal
	// app will receive these keypresses.
	newtios.Cc[syscall.VSUSP] = 032
	newtios.Cc[syscall.VINTR] = 0
	newtios.Cc[syscall.VQUIT] = 0

	if err = termios.Tcsetattr(fd, termios.TCSANOW, &newtios); err != nil {
		goto failed
	}

	t.set = true

	return nil

failed:
	if t.out != nil {
		_ = t.out.Close()
		t.out = nil
	}
	return err
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 78
// End:
