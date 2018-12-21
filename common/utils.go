package common

import (
	"fmt"
	"io"
	"os"
)

func ReadFileOctets(filePath string) *Octets {
	allData := BinRead(filePath)
	fs := Wrap(allData)
	return fs
}

func BinRead(filePath string) []byte {

	//path.Join("bin", "numbers.binary")
	fp, _ := os.Open(filePath)
	defer fp.Close()

	data := make([]byte, 16)
	var allData []byte
	//var k int32
	//var k []int8
	for {
		data = data[:cap(data)]

		// read bytes to slice
		n, err := fp.Read(data)
		//fmt.Println(" file read ", n, err, data)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		allData = append(allData, data...)

		// convert bytes to int32
		data = data[:n]
		//binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &k)
		//fmt.Println(k)
	}

	//fmt.Println(" all data ", len(allData), cap(allData), allData)
	return allData
}
