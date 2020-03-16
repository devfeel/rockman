package packet

type JsonRequest struct {
	Version string
	Command string
	Message interface{}
}
