package collector

import (
	"github.com/tkddnr924/gomft/bootsect"
	"github.com/tkddnr924/gomft/fragment"
	"github.com/tkddnr924/gomft/mft"

	"fmt"
	"io"
	"os"

	"syscall"
)

const supportedOemId = "NTFS    "

type Collector struct {
	Target      string // PATH
	Destination string // PATH
	VolumePath  string
	Volume      *os.File
}

func NewCollector(src string, dst string, letter string) *Collector {
	volume, err := os.Open(letter)

	if err != nil {
		fmt.Printf("Unable to open volume %q: %v\n", letter, err)
		return nil
	}

	return &Collector{
		Target:      src,
		Destination: dst,
		VolumePath:  letter,
		Volume:      volume,
	}
}

func (collector *Collector) CollectMFT() {
	// Parse Boot Sector
	bootSector := ParseBootSector(collector.Volume)

	if bootSector == nil {
		fmt.Println("Unable to parse boot sector")
		return
	}

	// Seek to $MFT
	bytesPerCluster := bootSector.BytesPerSector * bootSector.SectorsPerCluster
	mftPosInBytes := int64(bootSector.MftClusterNumber) * int64(bytesPerCluster)

	record := collector.SeekVolume(bootSector, mftPosInBytes)

	if record == nil {
		fmt.Println("Unable to seek to $MFT")
		return
	}

	dataRuns := collector.GetAttrData(record)

	// DataRun to Fragment
	fragments := mft.DataRunsToFragments(dataRuns, bytesPerCluster)
	totalLength := int64(0)

	for _, f := range fragments {
		totalLength += f.Length
	}

	// Copy $MFT
	outFile := collector.Destination + "\\$MFT"
	out, err := os.Create(outFile)

	if err != nil {
		fmt.Printf("Unable to create %q: %v\n", outFile, err)
		return
	}
	defer out.Close()

	length, err := PBCopy(out, fragment.NewReader(collector.Volume, fragments))

	if err != nil {
		fmt.Printf("Unable to copy $MFT: %v\n", err)
		return
	}

	if length != totalLength {
		fmt.Printf("Length mismatch: %d != %d\n", length, totalLength)
		return
	}

	info, _ := os.Stat(outFile)

	fmt.Printf("%+v\n", info.Sys())

	// attrs := _info.(map[string]interface{})

	// fmt.Printf("%+v\n", attrs)

	// if attrs, ok := _info.(map[string]interface{}); ok {
	// 	fmt.Println(getString(attrs, "created"))
	// } else {
	// 	fmt.Println("not ok", ok)
	// }
}

func getString(m map[string]interface{}, key string) string {
	if value, ok := m[key]; ok {
		if valueString, ok := value.(string); ok {
			return valueString
		}
	}
	return ""
}

// Get $FILE_NAME in MFT Record
func (collector *Collector) GetAttrFileName(record *mft.Record) *mft.FileName {
	fileNameAttribute := record.FindAttributes(mft.AttributeTypeFileName)
	if len(fileNameAttribute) == 0 {
		fmt.Println("No file name attributes found")
		return nil
	}

	fileName, err := mft.ParseFileName(fileNameAttribute[0].Data)
	if err != nil {
		fmt.Printf("Unable to parse file name: %v\n", err)
		return nil
	}

	return &fileName
}

// Get $DATA in MFT Record
func (collector *Collector) GetAttrData(record *mft.Record) []mft.DataRun {
	dataAttribute := record.FindAttributes(mft.AttributeTypeData)
	if len(dataAttribute) == 0 {
		fmt.Println("No data attributes found")
		return nil
	}

	dataRuns, err := mft.ParseDataRuns(dataAttribute[0].Data)
	if err != nil {
		fmt.Printf("Unable to parse data runs: %v\n", err)
		return nil
	}
	if len(dataRuns) == 0 {
		fmt.Println("No data runs found")
		return nil
	}

	return dataRuns
}

// Seek Volume
func (collector *Collector) SeekVolume(bootSector *bootsect.BootSector, position int64) *mft.Record {
	_, err := collector.Volume.Seek(position, 0)
	if err != nil {
		fmt.Printf("Unable to seek to %d: %v\n", position, err)
		return nil
	}

	mftSizeInBytes := bootSector.FileRecordSegmentSizeInBytes
	mftData := make([]byte, mftSizeInBytes)
	_, err = io.ReadFull(collector.Volume, mftData)
	if err != nil {
		fmt.Printf("Unable to read MFT data: %v\n", err)
		return nil
	}

	record, err := mft.ParseRecord(mftData)

	if err != nil {
		fmt.Printf("Unable to parse MFT record: %v\n", err)
		return nil
	}
	return &record
}

// Copy
func PBCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 1024*1024)
	return io.CopyBuffer(dst, src, buf)
}

// Parse Boot Sector
func ParseBootSector(volume *os.File) *bootsect.BootSector {
	src := (io.Reader)(volume)

	data := make([]byte, 512)
	_, err := io.ReadFull(src, data)

	if err != nil {
		fmt.Printf("Unable to read boot sector data: %v\n", err)
		return nil
	}

	bootSector, err := bootsect.Parse(data)

	if err != nil {
		fmt.Printf("Unable to parse boot sector data: %v\n", err)
		return nil
	}

	if bootSector.OemId != supportedOemId {
		fmt.Printf("Unknown OemId (file system type) %q (expected %q)\n", bootSector.OemId, supportedOemId)
		return nil
	}

	return &bootSector
}
