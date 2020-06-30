package isotp

import (
	"errors"
	"time"

	"github.com/x86ed/unCANnny/bits"
	"github.com/x86ed/unCANnny/can"
)

//PDU Persistent Data Unit or ISOTP FRAME
type PDU struct {
	Type      int
	Flag      int
	Length    int
	Data      []byte
	Raw       *can.Message
	Index     int
	BlockSize int
	SepTime   time.Duration
}

//GetType Returns the Current PDU/Frame type
func (f *PDU) GetType() (int, error) {
	typeInt := bits.GetNybble(f.Raw.Payload[0])[1]
	if typeInt < 0 || typeInt > 3 {
		f.Type = -1
		return f.Type, errors.New("Invalid ISOTP type")
	}
	f.Type = int(typeInt)
	return f.Type, nil
}

//Init creates a new PDU object (frame) from a set of bytes
func (f *PDU) Init(m can.Message) error {
	var err error
	f.Raw = &m
	_, err = f.GetType()
	switch f.Type {
	case FC:
		err = f.InitFlowControl()
	case CF:
		err = f.InitContinuationFrame()
	case FF:
		err = f.InitFirstFrame()
	case SF:
		err = f.InitSingleFrame()
	}
	return err
}

//InitFirstFrame creates the first frame
func (f *PDU) InitFirstFrame(b ...[]byte) error {
	var err error
	if len(b) > 0 {
		f.Raw.Payload = b[0]
		_, err = f.GetType()
		if err != nil {
			return err
		}
	}
	off, err := f.getFFDL()
	if err != nil {
		return err
	}
	f.Data = f.Raw.Payload[off:]
	return err
}

func (f *PDU) getFFDL() (int, error) {
	var err error
	minVal := 8
	if f.Raw.IsExtended {
		minVal = 7
	}
	if f.Length > 8 {
		minVal = len(f.Raw.Payload) - 1
		if f.Raw.IsExtended {
			minVal = len(f.Raw.Payload) - 2
		}

	}
	val, err := parseFFLength(f.Raw.Payload, minVal)
	if err != nil {
		return 0, err
	}
	if val < minVal {
		return 0, errors.New("length of first frame is too small")
	}
	return val, nil
}

func parseFFLength(b []byte, min int) (int, error) {
	n := bits.GetNybble(b[0])
	if n == 0 && b[1] == 0{
		
	}
}

//InitSingleFrame creates a single frame
func (f *PDU) InitSingleFrame(b ...[]byte) error {
	var err error
	if len(b) > 0 {
		f.Raw.Payload = b[0]
		_, err = f.GetType()
		if err != nil {
			return err
		}
	}
	err = f.getSFDL()
	if err != nil {
		return err
	}
	fl := len(f.Raw.Payload)
	off := fl - f.Length
	f.Data = f.Raw.Payload[off:]
	return err
}

func (f *PDU) getSFDL() error {
	upN := bits.GetNybble(f.Raw.Payload[0])[0]
	if upN != 0 {
		f.Length = int(upN)
		if f.Length > 7 || f.Length > len(f.Raw.Payload)-1 {
			return errors.New("invalid SF Data Length")
		}
		if f.Length == 7 && f.Raw.IsExtended {
			return errors.New("invalid length for extended frame")
		}
	}
	if len(f.Raw.Payload) < 2 {
		return errors.New("PDU length is longer than the data provided")
	}
	upN = f.Raw.Payload[1]
	f.Length = int(upN)
	if upN < 6 {
		return errors.New("invalid extended FF length")
	}
	if upN == 7 && !f.Raw.IsExtended {
		return errors.New("invalid length for non-extended frame")
	}
	if int(upN) > len(f.Raw.Payload)-2 {
		return errors.New("invalid SF Data Length")
	}
	return nil
}

//InitFlowControl creates a flow control frame
func (f *PDU) InitFlowControl(b ...[]byte) error {
	var err error
	if len(b) > 0 && len(b[0]) == 3 {
		f.Raw.Payload = b[0]
		_, err = f.GetType()
		if err != nil {
			return err
		}
	}
	if len(f.Raw.Payload) != 3 {
		return errors.New("Invalid Continuation Frame Length")
	}
	_, err = f.GetFlag()
	if err != nil {
		return err
	}
	f.BlockSize = int(f.Raw.Payload[1])
	_, err = f.GetSTMin()
	return err
}

//GetSTMin parses the ST_Min value for a FC frame
func (f *PDU) GetSTMin() (time.Duration, error) {
	stInt := int(f.Raw.Payload[2])
	if stInt >= 0 && stInt <= 0x7F {
		f.SepTime = time.Duration(stInt) * time.Millisecond
		return f.SepTime, nil
	} else if stInt >= 0xf1 && stInt <= 0xF9 {
		f.SepTime = time.Duration((stInt-0xF0)*100) * time.Microsecond
		return f.SepTime, nil
	}
	return 0, errors.New("Invalid ST Min value")
}

//GetFlag returns the flow control flag of the frame parsed from raw
func (f *PDU) GetFlag() (int, error) {
	flagInt := bits.GetNybble(f.Raw.Payload[0])[0]
	if flagInt < 0 || flagInt > 2 {
		f.Flag = -1
		return f.Flag, errors.New("Invalid Flow Control Flag")
	}
	f.Flag = int(flagInt)
	return f.Flag, nil
}

const (
	//SF = Single Frame type
	SF = iota
	//FF = First Frame type
	FF
	//CF = Continuation Frame type
	CF
	//FC = Flow Control Type
	FC
)

const (
	//Continue to send Flow Control flag
	Continue = iota
	//Wait flow control flag
	Wait
	//Overflow / Abort indicator flow control flag
	Overflow
)
