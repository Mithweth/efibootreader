//go:build darwin

package efi

import "io/fs"

type darwinBackend struct {
}

func (be *darwinBackend) IsEfi() bool {
	return true
}

func (be *darwinBackend) getVariablePath(name string) (string, error) {
	return "", fs.ErrExist
}

func newBackend() variableBackend {
	return &darwinBackend{}
}
