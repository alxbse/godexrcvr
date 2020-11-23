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
	} else if command == "generic" {
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
	} else if command == "rdpi" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := godexrcvr.OpenDevice(devicefile)
		defer device.Close()
		if err != nil {
			panic(err)
		}
		partInfo, err := godexrcvr.ReadDatabasePartionInfo(device)
		if err != nil {
			panic(err)
		}
		fmt.Printf("SchemaVersion:      %v\n", partInfo.SchemaVersion)
		fmt.Printf("PageHeaderVersion:  %v\n", partInfo.PageHeaderVersion)
		fmt.Printf("PageDataLength:     %v\n", partInfo.PageDataLength)
		for _, partition := range partInfo.Partitions {
			fmt.Printf("Partition:\n")
			fmt.Printf("  Name: %v\n", partition.Name)
			fmt.Printf("  Id: %v\n", partition.Id)
			fmt.Printf("  RecordRevision: %v\n", partition.RecordRevision)
			fmt.Printf("  RecordLength: %v\n", partition.RecordLength)
		}
	}
}
