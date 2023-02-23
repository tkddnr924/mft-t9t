package collector

import (
	// "github.com/tkddnr924/gomft/mft"
	"github.com/tkddnr924/gomft/bootsect"
	"io"
	"os"
)

const supportedOemId = "NTFS    "

type Collector struct {
	Target      string // PATH
	Destination string // PATH
	VolumePath  string
	Volume      *os.File
}

func NewCollector(src string, dst string, letter string) *Collector {

	volume, _ := os.Open(src)

	return &Collector{
		Target:      src,
		Destination: dst,
		VolumePath:  letter,
		Volume:      volume,
	}
}

func (collector *Collector) ParseBootSector() {
	data := make([]byte, 512)
	_, err := io.ReadFull(collector.Volume, data)

	if err != nil {
		return
	}

	bootSector, err := bootsect.Parse(data)

	if err != nil {
		fmt.Printf("Unable to parse boot sector data: %v\n", err)
	}

	if bootSector.OemId != supportedOemId {
		fmt.Printf("Unknown OemId (file system type) %q (expected %q)\n", bootSector.OemId, supportedOemId)
	}
}

// func (collector *Collector) Collect() {
// }
