package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func BinRead(filePath string) []byte {

	//path.Join("bin", "numbers.binary")
	fp, _ := os.Open(filePath)
	defer fp.Close()

	data := make([]byte, 4)
	var k int32
	for {
		data = data[:cap(data)]

		// read bytes to slice
		n, err := fp.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		// convert bytes to int32
		data = data[:n]
		binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &k)
		fmt.Println(k)
	}
	return nil
}
