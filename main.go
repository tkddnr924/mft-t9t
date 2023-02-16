package main

import (
	"fmt"
	"github.com/t9t/gomft/bootsect"
	"github.com/t9t/gomft/fragment"
	"github.com/t9t/gomft/mft"

	"io"
	"os"
)

func main() {
	volume := "\\\\.\\C:"

	in, err := os.Open(volume)
	if err != nil {
		fmt.Println("[1]", err)
	}
	defer in.Close()

	bootSectorData := make([]byte, 512)
	_, err = io.ReadFull(in, bootSectorData)
	if err != nil {
		fmt.Println("[2]", err)
	}
	defer in.Close()

	bootSector, err := bootsect.Parse(bootSectorData)

	if err != nil {
		fmt.Printf("[3] Unable to parse boot sector data: %v\n", err)
	}

	if bootSector.OemId != supportedOemId {
		fmt.Printf("[4] Unknown OemId (file system type) %q (expected %q)\n", bootSector.OemId, supportedOemId)
	}

	bytesPerCluster := bootSector.BytesPerSector * bootSector.SectorsPerCluster
	mftPosInBytes := int64(bootSector.MftClusterNumber) * int64(bytesPerCluster)

	_, err = in.Seek(mftPosInBytes, 0)
	if err != nil {
		fmt.Printf("[5] Unable to seek to MFT: %v\n", err)
	}

	mftSizeInBytes := bootSector.FileRecordSegmentSizeInBytes
	mftData := make([]byte, mftSizeInBytes)

	_, err = io.ReadFull(in, mftData)
	if err != nil {
		fmt.Printf("[6] Unable to read MFT: %v\n", err)
	}

	record, err := mft.ParseRecord(mftData)
	if err != nil {
		fmt.Printf("[7] Unable to parse MFT: %v\n", err)
	}

	dataAttributes := record.FindAttributes(mft.AttributeTypeData)
	if len(dataAttributes) == 0 {
		fmt.Println("[8] No data attributes found")
	}

	if len(dataAttributes) > 1 {
		fmt.Println("[9] More than one data attribute found")
	}

	dataAttribute := dataAttributes[0]
	if dataAttribute.Resident {
		fmt.Println("[10] Data attribute is resident")
	}

	dataRuns, err := mft.ParseDataRuns(dataAttribute.Data)
	if err != nil {
		fmt.Printf("[11] Unable to parse data runs: %v\n", err)
	}

	if len(dataRuns) == 0 {
		fmt.Println("[12] No data runs found")
	}

	fragments := mft.DataRunsToFragments(dataRuns, bytesPerCluster)
	totalLength := int64(0)

	for _, frag := range fragments {
		totalLength += int64(frag.Length)
	}

	outfile := "output/$MFT"
	out, err := openOutputFile(outfile)

	if err != nil {
		fmt.Println("[13]", err)
	}
	defer out.Close()

	n, err := CollectCopy(out, fragment.NewReader(in, fragments))

	if err != nil {
		fmt.Println("[14]", err)
	}

	if n != totalLength {
		fmt.Printf("[15] Expected to copy %d bytes, but copied only %d\\n", totalLength, n)
	}

	fmt.Println("END PROGRAM")
}

func openOutputFile(outfile string) (*os.File, error) {
	return os.Create(outfile)
}

func CollectCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 1024*1024)
	return io.CopyBuffer(dst, src, buf)
}

const supportedOemId = "NTFS    "
