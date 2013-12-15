package passwd

import "errors"

var passFunc func(prompt string) ([]byte, error)

// GetPasswd displays a prompt to, and reads in a password from, /dev/tty
// (/dev/cons). GetPasswd turns off character echoing while reading the
// password. The calling process should zero the password as soon as
// possible.
func GetPasswd(prompt string) ([]byte, error) {
	if passFunc == nil {
		return nil, errors.New("not supported")
	}
	return passFunc(prompt)
}
