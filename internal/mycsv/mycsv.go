package mycsv

import (
	"encoding/csv"
	"fmt"
	"os"
)

func WriteHrefCsv(href string, filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}

	w := csv.NewWriter(f)
	w.Write([]string{href})
	w.Flush()
}
