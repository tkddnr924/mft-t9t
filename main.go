package main

import (
	"fmt"
	"mft-t9t/collector"
	"os"
)

// func FindFileName(record mft.Record) (fileName mft.FileName) {
// 	attrs := record.FindAttributes(mft.AttributeTypeFileName)
// 	fileName, _ = mft.ParseFileName(attrs[0].Data)

// 	return fileName
// }

func main() {
	volume := "\\\\.\\C:"

	fmt.Println("START PROGRAM")

	if _, err := os.Stat("output"); os.IsNotExist(err) {
		os.Mkdir("output", 0777)
	}

	pb := collector.NewCollector("C:\\$MFT", "output", volume)

	pb.CollectMFT()

	fmt.Println("END PROGRAM")
}
