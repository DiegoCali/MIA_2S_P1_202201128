package interpreter

type MBR struct {
	Size       [4]byte
	TimeStamp  [4]byte
	Signature  [4]byte
	Fit        [1]byte
	Partitions [4]Partition
}

type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       [4]byte
	Size        [4]byte
	Name        [16]byte
	Correlative [4]byte
	Id          [4]byte
}
