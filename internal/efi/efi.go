package efi

import "os"

const EFIVARFS = "/sys/firmware/efi/efivars"

func IsEFI() bool {
	_, err := os.Stat(EFIVARFS)
	return err == nil
}
