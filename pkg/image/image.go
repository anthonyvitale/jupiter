// Package camera provides a way to take photos with a connected camera.
package image

import "io"

// Camera is the interface that provides a way to take pictures via a connected camera.
type Camera interface {
	Take() (io.Reader, error)
}

type PeripheralCamera struct{}
