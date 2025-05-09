package network

// Destination - a DICOM destination
type Destination struct {
	ID        string
	Name      string
	HostName  string
	CalledAE  string
	CallingAE string
	Port      int
	IsCStore  bool
	IsCFind   bool
	IsCMove   bool
	IsMWL     bool
	IsTLS     bool
	Anonymize bool
}
