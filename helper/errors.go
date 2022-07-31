package helper

import "fmt"

type errorInfo int8

const (
	PasswordLength errorInfo = iota
)

type ErrPasword struct {
	Info errorInfo
}

func (e ErrPasword) Error() string {
	switch e.Info {
	case PasswordLength:
		return fmt.Sprintf("pasword minimum length is %d", PasswordMinLength)
	default:
		return "password is not valid"
	}
}
