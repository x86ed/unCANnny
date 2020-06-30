package can

//Message describes a CAN message construct.
type Message struct {
	ID         uint32
	DLC        uint8
	Payload    []byte
	IsExtended bool
	IsFD       bool
}
