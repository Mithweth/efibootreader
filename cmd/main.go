package main

import (
	"fmt"
	"github.com/Mithweth/efibootreader/efi"
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

	bootIds, err := efi.GetBootIds()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var bootOrderIds []string
	bootOrderMapIds := make(map[uint16]struct{}, len(bootOrder))
	for _, id := range bootOrder {
		bootOrderIds = append(bootOrderIds, fmt.Sprintf("%04X", id))
		bootOrderMapIds[id] = struct{}{}
	}
	fmt.Printf("BootOrder: %s\n", strings.Join(bootOrderIds, ","))

	for _, id := range bootIds {
		bootEntry, err := efi.GetBootEntry(id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bootName := fmt.Sprintf("Boot%04X", id)
		if _, ok := bootOrderMapIds[id]; ok {
			bootName += "*"
		}
		fmt.Printf("%s %s %s\n", bootName, bootEntry.Description, bootEntry.DevicePath)
		//fmt.Printf("%#v\n", bootEntry.DevicePath)
	}
}
