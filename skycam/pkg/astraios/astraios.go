// Package astraios provides a way to take photos with a connected camera.
package astraios

import "io"

// Camera is the interface that provides a way to take pictures via a connected camera.
type Camera interface {
	TakePhoto() (io.Reader, error)
}

// ShellCamera provides a way to take a picture via the shell.
type ShellCamera interface {
	ExecuteCommand() error
}
