//go:build windows

package backend

import "io/fs"

type windowsBackend struct {
}

func (be *windowsBackend) IsEfi() bool {
	return true
}

func (be *windowsBackend) getVariablePath(name string) (string, error) {
	return "", fs.ErrExist
}

func NewBackend() variableBackend {
	return &windowsBackend{}
}
