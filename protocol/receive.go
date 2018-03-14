package protocol

type detector struct {
	Qcode      string
	ReceiverID string
	Disable    bool
}

type receiver struct {
	Qcode      string
	DetectorID string
	IPAddr     string
}
