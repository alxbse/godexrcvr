package main

import (
	"fmt"
	"os"

	"github.com/thecubic/godexrcvr"
	"github.com/urfave/cli/v2"
)

func connectAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}

	cmd := godexrcvr.CmdPing
	fmt.Printf("cmd: %s\n", cmd.String())
	defer device.Close()

	return nil
}

func pingAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	godexrcvr.DoAPing(device)

	return nil
}

func readFirmwareHeaderAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	godexrcvr.ReadFirmwareHeader(device)

	return nil
}

func readBatteryLevelAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	level, err := godexrcvr.ReadBatteryLevel(device)
	if err != nil {
		return err
	}

	fmt.Printf("battery level: %d\n", level)
	return nil
}

func readTransmitterIDAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	transmitterID, err := godexrcvr.ReadTransmitterID(device)
	if err != nil {
		return err
	}

	fmt.Printf("transmitter id: %s\n", transmitterID)
	return nil
}

func readDatabasePartitionInfoAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	databasePartionInfo, err := godexrcvr.ReadDatabasePartionInfo(device)
	if err != nil {
		return err
	}

	fmt.Printf("SchemaVersion: %s\n", databasePartionInfo.SchemaVersion)
	fmt.Printf("PageHeaderVersion: %s\n", databasePartionInfo.PageHeaderVersion)
	fmt.Printf("PageDataLength: %s\n", databasePartionInfo.PageDataLength)

	for _, partition := range databasePartionInfo.Partitions {
		fmt.Printf("Name: %s, ID: %s, RecordRevision: %s, RecordLength: %s \n", partition.Name, partition.Id, partition.RecordRevision, partition.RecordLength)
	}
	return nil
}

func readManufacturingDataAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	_, err = godexrcvr.ReadManufacturingData(device)
	if err != nil {
		return err
	}

	return nil
}

func readEgvDataAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	records, err := godexrcvr.ReadEgvData(device)
	if err != nil {
		return err
	}

	for _, record := range *records {
		mmol := MgToMmol(record.Value)
		fmt.Printf("glucose mg=%di,mmol=%di %d\n", record.Value, mmol, record.SystemTime.Unix())
	}

	return nil
}

func readGlucoseUnitAction(ctx *cli.Context) error {
	deviceFile := ctx.String("device-file")

	device, err := godexrcvr.OpenDevice(deviceFile)
	if err != nil {
		return err
	}
	defer device.Close()

	unit, err := godexrcvr.ReadGlucoseUnit(device)
	if err != nil {
		return err
	}

	fmt.Println(unit)
	return nil
}

func main() {
	app := cli.App{
		Name: "godexrcvr",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "device-file",
				Usage:    "device file",
				Required: true,
				EnvVars:  []string{"DEXCOM_DEVICE"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "connect",
				Action: connectAction,
			},
			{
				Name:   "ping",
				Usage:  "ping device",
				Action: pingAction,
			},
			{
				Name:   "read-firmware-header",
				Usage:  "read firmware header",
				Action: readFirmwareHeaderAction,
			},
			{
				Name:   "read-battery-level",
				Usage:  "read battery level",
				Action: readBatteryLevelAction,
			},
			{
				Name:   "read-transmitter-id",
				Usage:  "read transmitter id",
				Action: readTransmitterIDAction,
			},
			{
				Name:   "read-database-partition-info",
				Usage:  "read database partition info",
				Action: readDatabasePartitionInfoAction,
			},
			{
				Name:   "read-manufacturing-data",
				Usage:  "read manufacturing data",
				Action: readManufacturingDataAction,
			},
			{
				Name:   "read-egv-data",
				Usage:  "read egv data",
				Action: readEgvDataAction,
			},
			{
				Name:   "read-glucose-unit",
				Usage:  "read glucose unit",
				Action: readGlucoseUnitAction,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
