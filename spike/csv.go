package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func main() {

	reader := strings.NewReader("1,2,3\n3,2,1\n\n5,6,7")

	csvReader := csv.NewReader(reader)

	for {
		row, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(row)
	}

}
