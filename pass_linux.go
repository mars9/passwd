package passwd

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func ioctl(fd uintptr, cmd uintptr, data *syscall.Termios) error {
	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		cmd,
		uintptr(unsafe.Pointer(data)),
	); err != 0 {
		return syscall.ENOTTY
	}
	return nil
}

func GetPasswd(prompt string) ([]byte, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	defer tty.Close()

	var ostate syscall.Termios
	if err = ioctl(tty.Fd(), syscall.TCGETS, &ostate); err != nil {
		return nil, err
	}

	nstate := ostate
	nstate.Lflag &^= (syscall.ECHO | syscall.ISIG)
	if err = ioctl(tty.Fd(), syscall.TCSETS, &nstate); err != nil {
		return nil, err
	}
	defer ioctl(tty.Fd(), syscall.TCSETS, &ostate)

	fmt.Fprint(tty, prompt)
	r := bufio.NewReader(tty)
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	fmt.Fprint(tty, "\n")
	return line[:len(line)-1], nil
}
