//go:build windows

package efi

type windowsBackend struct {
	
}

func newBackend() variableBackend {
	return &windowsBackend{}
}

