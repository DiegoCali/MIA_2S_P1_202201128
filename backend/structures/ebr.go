package structures

// EBR Size in bytes: 33
type EBR struct {
	Mount int32    // bytes: 4
	Fit   [1]byte  // bytes: 1
	Start int32    // bytes: 4
	Size  int32    // bytes: 4
	Next  int32    // bytes: 4
	Name  [16]byte // bytes: 16
}
