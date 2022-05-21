package godexrcvr

import (
	"fmt"

	"github.com/google/gousb"
)

// SyncByte is a wire-signal for stuff
var SyncByte = 0x01

type DexcomCmd byte

// BASE_TIME = datetime.datetime(2009, 1, 1)
// DEXCOM_EPOCH = 1230768000

var (
	DexcomVendor        = gousb.ID(0x22a3)
	Gen4ReceiverProduct = gousb.ID(0x0047)
)

type DexcomPacket struct {
	cmd     DexcomCmd
	payload []byte
}

const (
	CmdNull                       DexcomCmd = 0x00
	CmdAck                        DexcomCmd = 0x01
	CmdNak                        DexcomCmd = 0x02
	CmdInvalidCommand             DexcomCmd = 0x03
	CmdInvalidParam               DexcomCmd = 0x04
	CmdIncompletePacketReceived   DexcomCmd = 0x05
	CmdReceiverError              DexcomCmd = 0x06
	CmdInvalidMode                DexcomCmd = 0x07
	CmdPing                       DexcomCmd = 0x0A
	CmdReadFirmwareHeader         DexcomCmd = 0x0B
	CmdReadDatabasePartitionInfo  DexcomCmd = 0x0F
	CmdReadDataPageRange          DexcomCmd = 0x10
	CmdReadDataPages              DexcomCmd = 0x11
	CmdReadDataPageHeader         DexcomCmd = 0x12
	CmdReadLanguage               DexcomCmd = 0x1B
	CmdReadDisplayTimeOffset      DexcomCmd = 0x1C
	CmdWriteDisplayTimeOffset     DexcomCmd = 0x1D
	CmdReadSystemTime             DexcomCmd = 0x22
	CmdReadSystemTimeOffset       DexcomCmd = 0x23
	CmdReadGlucoseUnit            DexcomCmd = 0x25
	CmdReadClockMode              DexcomCmd = 0x29
	CmdReadTransmitterID          DexcomCmd = 0x19
	CmdWriteTransmitterID         DexcomCmd = 0x1A
	CmdWriteLanguage              DexcomCmd = 0x1C
	CmdReadRTC                    DexcomCmd = 0x1F
	CmdResetReceiver              DexcomCmd = 0x20
	CmdReadBatteryLevel           DexcomCmd = 0x21
	CmdWriteSystemTime            DexcomCmd = 0x24
	CmdWriteGlucoseUnit           DexcomCmd = 0x26
	CmdReadBlindedMode            DexcomCmd = 0x27
	CmdWriteBlindedMode           DexcomCmd = 0x28
	CmdWriteClockMode             DexcomCmd = 0x2A
	CmdReadDeviceMode             DexcomCmd = 0x2B
	CmdEraseDatabase              DexcomCmd = 0x2D
	CmdShutdownReceiver           DexcomCmd = 0x2E
	CmdWritePcParameters          DexcomCmd = 0x2F
	CmdReadBatteryState           DexcomCmd = 0x30
	CmdReadHardwareBoardID        DexcomCmd = 0x31
	CmdReadFirmwareSettings       DexcomCmd = 0x36
	CmdReadEnableSetupWizardFlag  DexcomCmd = 0x37
	CmdReadSetupWizardState       DexcomCmd = 0x39
	CmdReadChargerCurrentSetting  DexcomCmd = 0x3b
	CmdWriteChargerCurrentSetting DexcomCmd = 0x3c
)

const (
	RecordTypeManufacturingData = 0x0
	RecordTypeEgvData           = 0x4
)

func (cmd DexcomCmd) String() string {
	switch cmd {
	case CmdNull:
		return "CmdNULL"
	case CmdPing:
		return "CmdPing"
	case CmdAck:
		return "CmdAck"
	case CmdNak:
		return "CmdNak"
	case CmdInvalidCommand:
		return "CmdInvalidCommand"
	case CmdInvalidParam:
		return "CmdInvalidParam"
	case CmdIncompletePacketReceived:
		return "CmdIncompletePacketReceived"
	case CmdReceiverError:
		return "CmdReceiverError"
	case CmdInvalidMode:
		return "CmdInvalidMode"

	case CmdReadFirmwareHeader:
		return "CmdReadFirmwareHeader"
	case CmdReadDatabasePartitionInfo:
		return "CmdReadDatabasePartitionInfo"
	case CmdReadDataPageRange:
		return "CmdReadDataPageRange"
	case CmdReadDataPages:
		return "CmdReadDataPages"
	case CmdReadDataPageHeader:
		return "CmdReadDataPageHeader"
	case CmdReadLanguage:
		return "CmdReadLanguage"
	case CmdReadDisplayTimeOffset:
		return "CmdReadDisplayTimeOffset"
	case CmdWriteDisplayTimeOffset:
		return "CmdWriteDisplayTimeOffset"
	case CmdReadSystemTime:
		return "CmdReadSystemTime"
	case CmdReadSystemTimeOffset:
		return "CmdReadSystemTimeOffset"
	case CmdReadGlucoseUnit:
		return "CmdReadGlucoseUnit"
	case CmdReadClockMode:
		return "CmdReadClockMode"
	}
	// TODO
	return "CmdUNKNOWN"
}

type FirmwareHeader struct {
	SchemaVersion      string `xml:"SchemaVersion,attr"`
	ApiVersion         string `xml:"ApiVersion,attr"`
	TestApiVersion     string `xml:"TestApiVersion,attr"`
	ProductId          string `xml:"ProductId,attr"`
	ProductName        string `xml:"ProductName,attr"`
	SoftwareNumber     string `xml:"SoftwareNumber,attr"`
	FirmwareVersion    string `xml:"FirmwareVersion,attr"`
	PortVersion        string `xml:"PortVersion,attr"`
	RFVersion          string `xml:"RFVersion,attr"`
	BLESoftwareVersion string `xml:"BLESoftwareVersion,attr"`
	BLEHardwareVersion string `xml:"BLEHardwareVersion,attr"`
	BLEDeviceAddress   string `xml:"BLEDeviceAddress,attr"`
	DexBootVersion     string `xml:"DexBootVersion,attr"`
}

type PartitionInfo struct {
	SchemaVersion     string      `xml:"SchemaVersion,attr"`
	PageHeaderVersion string      `xml:"PageHeaderVersion,attr"`
	PageDataLength    string      `xml:"PageDataLength,attr"`
	Partitions        []Partition `xml:"Partition"`
}

type Partition struct {
	Name           string `xml:"Name,attr"`
	Id             string `xml:"Id,attr"`
	RecordRevision string `xml:"RecordRevision,attr"`
	RecordLength   string `xml:"RecordLength,attr"`
}

type GlucoseUnit byte

const (
	GlucoseUnitMmolL = 0x1
	GlucoseUnitMgDl  = 0x2
)

func (u *GlucoseUnit) String() (string, error) {
	switch *u {
	case GlucoseUnitMmolL:
		return "mmol/l", nil
	case GlucoseUnitMgDl:
		return "mg/dl", nil

	}
	return "", fmt.Errorf("unknown glucose unit")
}
