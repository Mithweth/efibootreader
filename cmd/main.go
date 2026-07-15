package main

import (
	"fmt"
	"github.com/Mithweth/efibootreader/internal/efi"
	"os"
	"strings"
)

func main() {
	if !efi.IsEFI() {
		fmt.Println("Machine started in BIOS mode")
		os.Exit(1)
	}
	bootOrder, err := efi.GetBootOrder()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bootCurrent, err := efi.GetBootCurrent()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("BootCurrent: %04X\n", bootCurrent)
	bootOrderIds := make([]string, 0, len(bootOrder))
	for _, id := range bootOrder {
		bootOrderIds = append(bootOrderIds, fmt.Sprintf("%04X", id))
	}
	fmt.Printf("BootOrder: %s\n", strings.Join(bootOrderIds, ","))

	for _, bootId := range bootOrder {
		bootEntry, err := efi.GetBootEntry(bootId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Boot%04X: %s (%+v)\n", bootId, bootEntry.Description, bootEntry.DevicePath)

	}
}
