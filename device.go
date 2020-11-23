package godexrcvr

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"

	"github.com/google/gousb"

	// "go.bug.st/serial.v1"
	"time"

	"github.com/snksoft/crc"
	"github.com/tarm/serial"
)

var (
	DexcomVendor        = gousb.ID(0x22a3)
	Gen4ReceiverProduct = gousb.ID(0x0047)
)

type DexcomPacket struct {
	cmd     DexcomCmd
	payload []byte
}

// bs := make([]byte, 8)
// binary.PutUvarint(bs, uint64(64))

func DexcomFilter(desc *gousb.DeviceDesc) bool {
	switch desc.Vendor {
	case DexcomVendor:
		return desc.Product == Gen4ReceiverProduct
	default:
		return false
	}
}

func OpenDevice(device string) (*serial.Port, error) {
	sCfg := &serial.Config{Name: device, Baud: 115200, ReadTimeout: time.Second * 5}
	serDev, err := serial.OpenPort(sCfg)
	if err != nil {
		return nil, err
	}
	return serDev, nil
}

func checksumPacket(packet []byte) []byte {
	ckPacket := make([]byte, len(packet)+2)
	checksum := checksumForPacket(packet)
	copy(ckPacket, packet)
	binary.LittleEndian.PutUint16(ckPacket[len(packet):], checksum)
	return ckPacket
}

func checksumForPacket(packet []byte) uint16 {
	return uint16(crc.CalculateCRC(crc.XMODEM, packet))
}

func decodePacket(packet []byte) (*DexcomPacket, error) {
	if packet[0] != 0x01 {
		return nil, fmt.Errorf("packet does not start with sync byte")
	}
	pktFull := len(packet)
	packetLength := binary.LittleEndian.Uint16(packet[1:3])
	payloadLength := packetLength - 6
	packetChecksum := binary.LittleEndian.Uint16(packet[pktFull-2:])
	calcChecksum := checksumForPacket(packet[:len(packet)-2])
	if packetChecksum != calcChecksum {
		return nil, fmt.Errorf("checksum mismatch: %v [encoded], %v [calculated]", packetChecksum, calcChecksum)
	}
	payload := make([]byte, payloadLength)
	if payloadLength > 0 {
		copy(payload, packet[4:packetLength-2])
	}
	decodedPacket := &DexcomPacket{DexcomCmd(packet[3]), payload}
	return decodedPacket, nil
}

func buildCmdPacket(cmd DexcomCmd) []byte {
	var packet []byte
	packet = make([]byte, 6)
	packet[0] = 0x01
	binary.LittleEndian.PutUint16(packet[1:3], 6)
	packet[3] = byte(cmd)
	checksum := checksumForPacket(packet[:4])
	binary.LittleEndian.PutUint16(packet[4:6], checksum)
	return packet
}

func buildPacket(cmd DexcomCmd, payload []byte) []byte {
	var packet []byte
	payloadLength := uint16(len(payload))
	packetLength := payloadLength + 6
	packet = make([]byte, payloadLength+6)
	packet[0] = 0x01
	binary.LittleEndian.PutUint16(packet[1:3], packetLength)
	packet[3] = byte(cmd)
	checksum := checksumForPacket(packet[:packetLength-2])
	binary.LittleEndian.PutUint16(packet[packetLength-2:], checksum)
	return packet
}

func ReadPacket(device *serial.Port) (*DexcomPacket, error) {
	// this device is picky to read from.
	// instead of "give me what you got",
	// this will read the first 4 bytes
	// which is syncbyte packetlen cmd
	// and then tape-read the rest
	var err error
	var n int
	phdr := make([]byte, 4)
	n, err = device.Read(phdr)
	if err != nil {
		return nil, err
	}
	if n < 4 {
		return nil, fmt.Errorf("underflow")
	}
	if phdr[0] != 0x01 {
		return nil, fmt.Errorf("packet does not start with sync byte")
	}
	packetLength := binary.LittleEndian.Uint16(phdr[1:3])

	remaining := int(packetLength - uint16(n))
	toread := remaining

	fullPacket := make([]byte, 4+remaining)
	copy(fullPacket[0:4], phdr)

	rbuf := make([]byte, 128)
	offset := 4
	for toread > 0 {
		n, err = device.Read(rbuf)
		if err != nil {
			return nil, err
		}
		copy(fullPacket[offset:], rbuf)
		toread = toread - n
		offset = offset + n
	}

	return decodePacket(fullPacket)
}

func ReadBatteryLevel(device *serial.Port) (int, error) {
	packet := buildCmdPacket(CmdReadBatteryLevel)
	_, err := device.Write(packet)
	if err != nil {
		return -1, err
	}
	device.Flush()

	pkt, err := ReadPacket(device)
	if err != nil {
		return -1, err
	}
	if pkt.cmd != CmdAck {
		return -1, fmt.Errorf("command not acked")
	}
	return int(pkt.payload[0]), nil
}

func ReadFirmwareHeader(device *serial.Port) {
	packet := buildCmdPacket(CmdReadFirmwareHeader)
	_, err := device.Write(packet)
	if err != nil {
		fmt.Printf("error in writing: %v\n", err)
	}
	device.Flush()

	pkt, err := ReadPacket(device)
	if err != nil {
		panic(err)
	}
	if pkt.cmd != CmdAck {
		panic("command not acked")
	}
	fwHdr := new(FirmwareHeader)
	err = xml.Unmarshal(pkt.payload, &fwHdr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("SchemaVersion:      %v\n", fwHdr.SchemaVersion)
	fmt.Printf("ApiVersion:         %v\n", fwHdr.ApiVersion)
	fmt.Printf("TestApiVersion:     %v\n", fwHdr.TestApiVersion)
	fmt.Printf("ProductId:          %v\n", fwHdr.ProductId)
	fmt.Printf("ProductName:        %v\n", fwHdr.ProductName)
	fmt.Printf("SoftwareNumber:     %v\n", fwHdr.SoftwareNumber)
	fmt.Printf("FirmwareVersion:    %v\n", fwHdr.FirmwareVersion)
	fmt.Printf("PortVersion:        %v\n", fwHdr.PortVersion)
	fmt.Printf("RFVersion:          %v\n", fwHdr.RFVersion)
	fmt.Printf("BLESoftwareVersion: %v\n", fwHdr.BLESoftwareVersion)
	fmt.Printf("BLEHardwareVersion: %v\n", fwHdr.BLEHardwareVersion)
	fmt.Printf("BLEDeviceAddress:   %v\n", fwHdr.BLEDeviceAddress)
	fmt.Printf("DexBootVersion:     %v\n", fwHdr.DexBootVersion)
}

func DoAPing(device *serial.Port) bool {
	var err error
	var n int
	packet := buildCmdPacket(CmdPing)

	n, err = device.Write(packet)
	if err != nil {
		fmt.Printf("error in writing: %v\n", err)
	} else {
		fmt.Printf("nW: %v\n", n)
	}
	device.Flush()

	rbuf := make([]byte, 2560)
	n, err = device.Read(rbuf)
	if err != nil {
		fmt.Printf("error in reading: %v\n", err)
	} else {
		fmt.Printf("nR: %v\n", n)
		fmt.Printf("packet: %v\n", rbuf[:n])
	}
	ackPacket, err := decodePacket(rbuf[:n])
	return ackPacket.cmd == CmdAck
}

func ReadGeneric(device *serial.Port, cmd DexcomCmd) error {
	packet := buildCmdPacket(cmd)
	_, err := device.Write(packet)
	if err != nil {
		return err
	}
	device.Flush()

	pkt, err := ReadPacket(device)
	if err != nil {
		return err
	}
	if pkt.cmd != CmdAck {
		return fmt.Errorf("command not acked")
	}
	fmt.Printf("return payload (byte): %v\n", pkt.payload)
	fmt.Printf("return payload (str): %v\n", string(pkt.payload))
	return nil
}

func ReadTransmitterID(device *serial.Port) (string, error) {
	packet := buildCmdPacket(CmdReadTransmitterID)
	_, err := device.Write(packet)
	if err != nil {
		return "", err
	}
	device.Flush()

	pkt, err := ReadPacket(device)
	if err != nil {
		return "", err
	}
	if pkt.cmd != CmdAck {
		return "", fmt.Errorf("command not acked")
	}
	return string(pkt.payload), nil
}

//
// )

// func main() {
//         c := &serial.Config{Name: "COM45", Baud: 115200}
//         s, err := serial.OpenPort(c)
//         if err != nil {
//                 log.Fatal(err)
//         }

//         n, err := s.Write([]byte("test"))
//         if err != nil {
//                 log.Fatal(err)
//         }

//         buf := make([]byte, 128)
//         n, err = s.Read(buf)
//         if err != nil {
//                 log.Fatal(err)
//         }
// 		log.Printf("%q", buf[:n])

// module.exports.calcCRC = function (bytes, size, initial_remainder, final_xor) {
// 	var crc16;
// 	var i, j;

// 	crc16 = initial_remainder;
// 	// Divide the buffer by the polynomial, a byte at a time.
// 	for (i = 0; i < size; i++) {
// 	  crc16 = this.CRC_TABLE[(bytes[i] ^ (crc16 >> 8)) & 0xFF] ^ ((crc16 << 8) & 0xFFFF);
// 	}
// 	// The final remainder is the CRC.
// 	return (crc16 ^ final_xor);
//   };

// func CalcCRC(bytes []byte) uint16 {
// 	// initial remainder
// 	var crc16 = 0x0000

// }
// // if (this.testCRC_D('\x02\x06\x06\x03') != 50445) {
//     console.log('CRC_D logic is NOT CORRECT!!!');
//     return false;
//   }
//   return true;

// module.exports.calcCRC_D = function (bytes, size) {
// 	return this.calcCRC(bytes, size,
// 								 this.D_INITIAL_REMAINDER,
// 								 this.D_FINAL_XOR_VALUE);
//   };
// module.exports.D_INITIAL_REMAINDER = 0x0000;
// module.exports.D_FINAL_XOR_VALUE = 0x0000;
// module.exports.CRC_TABLE = [
//   0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
//   0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
//   0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
//   0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
//   0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
//   0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
//   0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
//   0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
//   0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
//   0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
//   0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
//   0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
//   0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
//   0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
//   0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
//   0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
//   0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
//   0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
//   0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
//   0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
//   0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
//   0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
//   0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
//   0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
//   0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
//   0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
//   0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
//   0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
//   0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
//   0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
//   0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
//   0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0
// ];

// var buildPacket = function (command, payloadLength, payload) {
//     var datalen = payloadLength + 6;
//     var buf = new ArrayBuffer(datalen);
//     var bytes = new Uint8Array(buf);
//     var ctr = struct.pack(bytes, 0, 'bsb', SYNC_BYTE,
//                           datalen, command);
//     ctr += struct.copyBytes(bytes, ctr, payload, payloadLength);
//     var crc = crcCalculator.calcCRC_D(bytes, ctr);
//     struct.pack(bytes, ctr, 's', crc);
//     return buf;
//   };

//   var readFirmwareHeader = function () {
//     return {
//       packet: buildPacket(
//         CMDS.READ_FIRMWARE_HEADER.value, 0, null
//       ),
//       parser: function (packet) {
//         var data = parseXMLPayload(packet);
//         firmwareHeader = data;
//         return data;
//       }
//     };
//   };

// BASE_TIME = datetime.datetime(2009, 1, 1)
// DEXCOM_EPOCH = 1230768000

// NULL = 0
// ACK = 1
// NAK = 2
// INVALID_COMMAND = 3
// INVALID_PARAM = 4
// INCOMPLETE_PACKET_RECEIVED = 5
// RECEIVER_ERROR = 6
// INVALID_MODE = 7
// PING = 10
// READ_FIRMWARE_HEADER = 11
// READ_DATABASE_PARTITION_INFO = 15
// READ_DATABASE_PAGE_RANGE = 16
// READ_DATABASE_PAGES = 17
// READ_DATABASE_PAGE_HEADER = 18
// READ_TRANSMITTER_ID = 25
// WRITE_TRANSMITTER_ID = 26
// READ_LANGUAGE = 27
// WRITE_LANGUAGE = 28
// READ_DISPLAY_TIME_OFFSET = 29
// WRITE_DISPLAY_TIME_OFFSET = 30
// READ_RTC = 31
// RESET_RECEIVER = 32
// READ_BATTERY_LEVEL = 33
// READ_SYSTEM_TIME = 34
// READ_SYSTEM_TIME_OFFSET = 35
// WRITE_SYSTEM_TIME = 36
// READ_GLUCOSE_UNIT = 37
// WRITE_GLUCOSE_UNIT = 38
// READ_BLINDED_MODE = 39
// WRITE_BLINDED_MODE = 40
// READ_CLOCK_MODE = 41
// WRITE_CLOCK_MODE = 42
// READ_DEVICE_MODE = 43
// ERASE_DATABASE = 45
// SHUTDOWN_RECEIVER = 46
// WRITE_PC_PARAMETERS = 47
// READ_BATTERY_STATE = 48
// READ_HARDWARE_BOARD_ID = 49
// READ_FIRMWARE_SETTINGS = 54
// READ_ENABLE_SETUP_WIZARD_FLAG = 55
// READ_SETUP_WIZARD_STATE = 57
// READ_CHARGER_CURRENT_SETTING = 59
// WRITE_CHARGER_CURRENT_SETTING = 60
// MAX_COMMAND = 61
// MAX_POSSIBLE_COMMAND = 255
