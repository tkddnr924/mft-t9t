package main

import (
	"fmt"
	"github.com/tkddnr924/gomft/bootsect"
	// "github.com/tkddnr924/gomft/fragment"
	"github.com/tkddnr924/gomft/mft"

	"mft-t9t/collector"

	"io"
	"os"
)

const supportedOemId = "NTFS    "

func ParseBootSector(src io.Reader) *bootsect.BootSector {
	data := make([]byte, 512)
	_, err := io.ReadFull(src, data)
	if err != nil {
		return nil
	}

	bootSector, err := bootsect.Parse(data)

	if err != nil {
		fmt.Printf("[3] Unable to parse boot sector data: %v\n", err)
	}

	if bootSector.OemId != supportedOemId {
		fmt.Printf("[4] Unknown OemId (file system type) %q (expected %q)\n", bootSector.OemId, supportedOemId)
	}

	return &bootSector
}

func FindFileName(record mft.Record) (fileName mft.FileName) {
	attrs := record.FindAttributes(mft.AttributeTypeFileName)
	fileName, _ = mft.ParseFileName(attrs[0].Data)

	return fileName
}

func main() {
	volume := "\\\\.\\C:"

	fmt.Println("START PROGRAM")

	if _, err := os.Stat("output"); os.IsNotExist(err) {
		os.Mkdir("output", 0777)
	}

	collector := collector.NewCollector("C:\\$MFT", "output", volume)

	fmt.Println(collector)

	// fmt.Println("[*] Open Volume", volume)
	// in, err := os.Open(volume)
	// if err != nil {
	// 	fmt.Println("[1]", err)
	// }
	// defer in.Close()

	// fmt.Println("[*] Read Boot Sector")
	// bootSector := ParseBootSector(in)

	// fmt.Println("[*] Read MFT & Seek $MFT")
	// bytesPerCluster := bootSector.BytesPerSector * bootSector.SectorsPerCluster
	// mftPosInBytes := int64(bootSector.MftClusterNumber) * int64(bytesPerCluster)

	// _, err = in.Seek(mftPosInBytes, 0)
	// if err != nil {
	// 	fmt.Printf("[5] Unable to seek to MFT: %v\n", err)
	// }

	// fmt.Println("[*] Parse MFT DATA")
	// mftSizeInBytes := bootSector.FileRecordSegmentSizeInBytes // 1024 bytes
	// mftData := make([]byte, mftSizeInBytes)

	// _, err = io.ReadFull(in, mftData)
	// if err != nil {
	// 	fmt.Printf("[6] Unable to read MFT: %v\n", err)
	// }

	// fmt.Println("[*] Parse MFT Record")
	// record, err := mft.ParseRecord(mftData)
	// if err != nil {
	// 	fmt.Printf("[7] Unable to parse MFT: %v\n", err)
	// }

	// fileName := FindFileName(record)

	// fmt.Println(fileName)

	// fmt.Println("[*] Parse Data Attribute")
	// dataAttributes := record.FindAttributes(mft.AttributeTypeData)
	// if len(dataAttributes) == 0 {
	// 	fmt.Println("[8] No data attributes found")
	// }

	// if len(dataAttributes) > 1 {
	// 	fmt.Println("[9] More than one data attribute found")
	// }

	// dataAttribute := dataAttributes[0]
	// if dataAttribute.Resident {
	// 	fmt.Println("[10] Data attribute is resident")
	// }

	// fmt.Println("[*] Parse Data Runs")
	// dataRuns, err := mft.ParseDataRuns(dataAttribute.Data)

	// if err != nil {
	// 	fmt.Printf("[11] Unable to parse data runs: %v\n", err)
	// }

	// if len(dataRuns) == 0 {
	// 	fmt.Println("[12] No data runs found")
	// }

	// fmt.Println("[*] Parse Data Fragments")
	// fragments := mft.DataRunsToFragments(dataRuns, bytesPerCluster)
	// totalLength := int64(0)

	// for _, frag := range fragments {
	// 	totalLength += int64(frag.Length)
	// }

	// outfile := "output/$MFT"
	// out, err := openOutputFile(outfile)

	// if err != nil {
	// 	fmt.Println("[13]", err)
	// }
	// defer out.Close()

	// fmt.Println("[*] Copy Data")
	// n, err := CollectCopy(out, fragment.NewReader(in, fragments))

	// if err != nil {
	// 	fmt.Println("[14]", err)
	// }

	// if n != totalLength {
	// 	fmt.Printf("[15] Expected to copy %d bytes, but copied only %d\\n", totalLength, n)
	// }

	fmt.Println("END PROGRAM")
}

func openOutputFile(outfile string) (*os.File, error) {
	return os.Create(outfile)
}

func CollectCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 1024*1024)
	return io.CopyBuffer(dst, src, buf)
}
