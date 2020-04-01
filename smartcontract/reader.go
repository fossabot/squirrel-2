package smartcontract

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type scriptContext struct {
	Position uint64
	Context  []byte
}

// OpCodeData stores opcode with its data.
type OpCodeData struct {
	OpCode byte
	Data   []byte
}

// DataStack is the stack of opcodes with its data(opcode without data is omitted).
// Stack direction: bottom(0) ---> top(N)
type DataStack []*OpCodeData

// Copy returns a deep copied data stack.
func (d *DataStack) Copy() *DataStack {
	newD := DataStack{}
	for _, opCodeData := range *d {
		newD.push(opCodeData.OpCode, opCodeData.Data)
	}

	return &newD
}

func (d *DataStack) push(opCode byte, data []byte) {
	item := &OpCodeData{opCode, data}
	*d = append(*d, item)
}

// PopData pops item data value from top.
func (d *DataStack) PopData() []byte {
	_, data := (*d).PopItem()
	return data
}

// PopItem pops item with value from top.
func (d *DataStack) PopItem() (byte, []byte) {
	if len(*d) == 0 {
		panic("Can not pop item from stack. Stack is empty.")
	}

	item := (*d)[len(*d)-1]
	*d = (*d)[:len(*d)-1]
	return item.OpCode, item.Data
}

func (sc *scriptContext) readNextOpCode() (byte, bool) {
	if sc.Position+1 >= uint64(len(sc.Context)) {
		return 0, false
	}
	opCode := sc.Context[sc.Position]
	sc.Position++
	return opCode, true
}

func (sc *scriptContext) readByte() byte {
	data := sc.Context[sc.Position]
	sc.Position++
	return data
}

func (sc *scriptContext) readUint16() uint16 {
	data := sc.Context[sc.Position : sc.Position+2]
	sc.Position += 2
	return binary.LittleEndian.Uint16(data)
}

func (sc *scriptContext) readUint32() uint32 {
	data := sc.Context[sc.Position : sc.Position+4]
	sc.Position += 4
	return binary.LittleEndian.Uint32(data)
}

func (sc *scriptContext) readUint64() uint64 {
	data := sc.Context[sc.Position : sc.Position+8]
	sc.Position += 8
	return binary.LittleEndian.Uint64(data)
}

func (sc *scriptContext) readBytes(byteLength uint64) ([]byte, error) {
	if sc.Position+byteLength > uint64(len(sc.Context)) {
		return nil, fmt.Errorf("failed to read bytes")
	}
	data := sc.Context[sc.Position : sc.Position+byteLength]
	sc.Position += byteLength
	return data, nil
}

func (sc *scriptContext) readVarBytes() ([]byte, error) {
	length := sc.readVarInt()
	data, err := sc.readBytes(length)
	return data, err
}

func (sc *scriptContext) readVarInt() uint64 {
	prefix := sc.readByte()
	if prefix == 0xFD {
		return uint64(sc.readUint16())
	} else if prefix == 0xFE {
		return uint64(sc.readUint32())
	} else if prefix == 0xFF {
		return sc.readUint64()
	}
	return uint64(prefix)
}

// ReadScript reads script string and extract stack of opcode with data.
func ReadScript(script string) *DataStack {
	bytes, _ := hex.DecodeString(script)
	context := scriptContext{0, bytes}
	stack := DataStack{}

	for {
		opCode, ok := context.readNextOpCode()
		if !ok {
			break
		}

		if opCode == 0x66 {
			return &stack
		}

		data, err := getDataFromScript(opCode, &context)
		if err != nil {
			return &stack
		}

		if data != nil {
			stack.push(opCode, data)
		}
	}

	return &stack
}

func getDataFromScript(opCode byte, context *scriptContext) ([]byte, error) {
	switch {
	case opCode == 0x00:
		return []byte{0}, nil
	case opCode >= 0x01 && opCode <= 0x4B:
		data, err := context.readBytes(uint64(opCode))
		return data, err
	case opCode == 0x4C:
		dataLen := context.readByte()
		data, err := context.readBytes(uint64(dataLen))
		return data, err
	case opCode == 0x4D:
		dataLen := context.readUint16()
		data, err := context.readBytes(uint64(dataLen))
		return data, err
	case opCode == 0x4E:
		dataLen := context.readUint32()
		data, err := context.readBytes(uint64(dataLen))
		return data, err
	case opCode == 0x4F:
		return []byte{0xFF, 0xFF}, nil
	case opCode >= 0x51 && opCode <= 0x60:
		data := opCode - 0x51 + 1
		return []byte{data}, nil
	case opCode == 0x61:
		return nil, nil
	case opCode >= 0x62 && opCode <= 0x64:
		value := context.readUint16()
		data := make([]byte, 2)
		binary.LittleEndian.PutUint16(data, value)
		return data, nil
	case opCode == 0x65:
		data, err := context.readBytes(2)
		return data, err
	case opCode == 0x66:
		return nil, nil
	case opCode == 0x67:
		data, err := context.readBytes(20)
		return data, err
	case opCode == 0x68:
		data, err := context.readVarBytes()
		return data, err
	case opCode == 0x69:
		data, err := context.readBytes(20)
		return data, err
	case opCode >= 0x6A && opCode <= 0x6D:
		return nil, nil
	case opCode >= 0x72 && opCode <= 0x7F:
		return nil, nil
	case opCode >= 0x80 && opCode <= 0x87:
		return nil, nil
	case opCode >= 0x8B && opCode <= 0x8D:
		return nil, nil
	case opCode == 0x8F:
		return nil, nil
	case opCode >= 0x90 && opCode <= 0x9C:
		return nil, nil
	case opCode >= 0x9E && opCode <= 0x9F:
		return nil, nil
	case opCode >= 0xA0 && opCode <= 0xAA:
		return nil, nil
	case opCode == 0xAC || opCode == 0xAE:
		return nil, nil
	case opCode >= 0xC0 && opCode <= 0xCD:
		return nil, nil
	case opCode >= 0xF0 && opCode <= 0xF1:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported opCode: %#02x", opCode)
	}
}
