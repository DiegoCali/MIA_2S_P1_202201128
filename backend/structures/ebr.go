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

func (ebr *EBR) Set(mount int32, fit string, start int32, size int32, next int32, name string) error {
	ebr.Mount = mount
	copy(ebr.Fit[:], fit)
	ebr.Start = start
	ebr.Size = size
	ebr.Next = next
	copy(ebr.Name[:], name)
	return nil
}
