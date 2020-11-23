package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/thecubic/godexrcvr"
)

var (
	debug = flag.Bool("debug", false, "enable debugging messages")
)

func usageQuit() {
	fmt.Println("usage: gdr-info <command> [args]")
	fmt.Println("                connect deviceFile")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	command := flag.Arg(0)
	if command == "" {
		usageQuit()
	} else if command == "hello" {
		fmt.Printf("syncbyte: %v\n", godexrcvr.SyncByte)
	} else if command == "connect" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		if err != nil {
			panic(err)
		}

		cmd := godexrcvr.CmdPing
		fmt.Printf("cmd: %v\n", cmd.String())
		defer device.Close()
	} else if command == "ping" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		if err != nil {
			panic(err)
		}

		godexrcvr.DoAPing(device)
		fmt.Println("haha! yes!")
		defer device.Close()
	} else if command == "fwhdr" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		defer device.Close()
		if err != nil {
			panic(err)
		}
		godexrcvr.ReadFirmwareHeader(device)
	} else if command == "rbat" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		defer device.Close()
		if err != nil {
			panic(err)
		}
		level, err := godexrcvr.ReadBatteryLevel(device)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Battery Level: %v%%\n", level)
	} else if command == "rxmtrid" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		defer device.Close()
		if err != nil {
			panic(err)
		}
		tid, err := godexrcvr.ReadTransmitterID(device)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transmitter ID: %v\n", tid)
	} else if command == "rxmtrid" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		defer device.Close()
		if err != nil {
			panic(err)
		}
		err = godexrcvr.ReadGeneric(device, godexrcvr.CmdReadTransmitterID)
		if err != nil {
			panic(err)
		}
	}

	// else if command == "crc" {
	// 	test := []byte{0x02, 0x06, 0x06, 0x03}

	// 	fmt.Printf("X23: %v\n", crc.CalculateCRC(crc.X25, test))
	// 	fmt.Printf("CCITT: %v\n", crc.CalculateCRC(crc.CCITT, test))
	// 	fmt.Printf("CRC16: %v\n", crc.CalculateCRC(crc.CRC16, test))
	// 	fmt.Printf("XMODEM: %v\n", crc.CalculateCRC(crc.XMODEM, test))
	// 	fmt.Printf("XMODEM2: %v\n", crc.CalculateCRC(crc.XMODEM2, test))

	// 	// fmt.Printf("CCITT: %v\n", crc16.ChecksumCCITT(test))
	// 	// fmt.Printf("CCITTFalse: %v\n", crc16.ChecksumCCITTFalse(test))
	// 	// fmt.Printf("IBM: %v\n", crc16.ChecksumIBM(test))
	// 	// fmt.Printf("MBus: %v\n", crc16.ChecksumMBus(test))
	// 	// fmt.Printf("SCSI: %v\n", crc16.ChecksumSCSI(test))
	// }
}
