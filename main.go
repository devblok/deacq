package main

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func main() {
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bufferedFile := bufio.NewReader(f)

	csv, err := os.Create(filename + ".csv.gz")
	if err != nil {
		panic(err)
	}
	defer csv.Close()

	gzipOutput := gzip.NewWriter(csv)
	defer gzipOutput.Close()

	bufferedOutput := bufio.NewWriter(gzipOutput)
	defer bufferedOutput.Flush()

	printCsvHeader(bufferedOutput, "PMT1 [volts]", "PMT2 [volts]", "PMT3 [volts]", "PMT4 [volts]", "Sorting pulse [volts]")

	result := make([]int32, 5)
	for {
		if err := binary.Read(bufferedFile, binary.BigEndian, &result); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		printCsv(bufferedOutput, result...)
	}
}

func printCsvHeader(w io.Writer, headers ...string) {
	for _, s := range headers {
		fmt.Fprintf(w, "%s\t", s)
	}
	fmt.Fprintf(w, "\n")
}

func printCsv(w io.Writer, values ...int32) {
	for _, v := range values {
		fmt.Fprintf(w, "%.2f\t", float32(v)/3276.8)
	}
	fmt.Fprintf(w, "\n")
}
