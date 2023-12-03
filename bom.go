package main

import (
	"fmt"
	"io"
	"os"
	//"io/ioutil"
)

const (
	bom   = "\uFEFF"
	usage = `
Usage: bom [options] <file>

options:
      -h, --help       Show this help message and exit.

      -q, --quiet      Do not show any hint.

      -d, --detect     Check if a UTF-8 BOM is at the start of the file.

      -a, --add        Add a UTF-8 BOM to the start of the file if such a BOM is not found. This does not check if the file is UTF-8 encoded or not.

      -s, --strip      Remove the UTF-8 BOM from the start of the file if such a BOM is found. This does not check if the file is UTF-8 encoded or not.

By default, this program checks if a UTF-8 BOM is at the start of a file. If one is found then it's removed, otherwise a UTF-8 BOM will be added to the start of the file. It does not check if the file is UTF-8 encoded or not.
`
)

func detectBOM(bs []byte) bool {
	const n = len(bom) //3
	if len(bs) < n {
		return false
	}
	for i := 0; i < n; i++ {
		if bs[i] != bom[i] {
			return false
		}
	}
	return true
}

func main() {
	var (
		quiet                       = false
		shouldStripIfAny            = false
		shouldAddIfNone             = false
		shouldOutputDetectionResult = false
		endOfFileReached            = false
		fileName                    string
	)
	if len(os.Args) == 1 {
		fmt.Print(usage)
		return
	}
	for _, a := range os.Args[1:] {
		switch a {
		case "-h",
			"--help":
			fmt.Print(usage)
			return
		case "-q",
			"--quiet":
			quiet = true
		case "-d",
			"--detect":
			shouldOutputDetectionResult = true
		case "-a",
			"--add":
			shouldAddIfNone = true
		case "-s",
			"--strip":
			shouldStripIfAny = true
		default:
			fileName = a
		}
	}
	if !(shouldAddIfNone || shouldStripIfAny || shouldOutputDetectionResult) {
		shouldAddIfNone = true
		shouldStripIfAny = true
	}
	inputFile, err := os.Open(fileName)
	if err != nil {
	    fmt.Println(err.Error())
		os.Exit(1)
	}
	defer inputFile.Close() // Error if already closed, but it'll be ignored.
	buffer := make([]byte, len(bom))
	n, err := inputFile.Read(buffer)
	if err != nil {
		if err == io.EOF {
			endOfFileReached = true
		} else {
			panic(err)
		}
	}
	if detectBOM(buffer[:n]) {
		if shouldOutputDetectionResult && !quiet {
			fmt.Println("BOM found.")
		}
		if shouldStripIfAny {
			tmp, err := os.CreateTemp("", "") //ioutil.Tempfile("", "")
			if err != nil {
				panic(err)
			}
			defer os.Remove(tmp.Name())
			defer tmp.Close()
			_, err = io.Copy(tmp, inputFile)
			if err != nil {
				panic(err)
			}
			st, err := inputFile.Stat()
			if err != nil {
				panic(err)
			}
			mode := st.Mode()
			err = inputFile.Close()
			if err != nil {
				panic(err)
			}
			outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, mode) // Overwrite the file but preserve its mode.
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()
			_, err = tmp.Seek(0, 0)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(outputFile, tmp)
			if err != nil {
				panic(err)
			}
			if !quiet {
				fmt.Println("BOM removed.")
			}
		} else if shouldAddIfNone && !quiet {
		    if !shouldOutputDetectionResult {
		        //The information is still displayed when it's a no-op.
		        fmt.Println("BOM found.")
		    }
		    fmt.Println("Nothing is added.")
		}
	} else {
		if shouldOutputDetectionResult && !quiet {
			fmt.Println("BOM not found.")
		}
		if shouldAddIfNone {
			tmp, err := os.CreateTemp("", "") //ioutil.Tempfile("", "")
			if err != nil {
				panic(err)
			}
			defer os.Remove(tmp.Name())
			defer tmp.Close()

			_, err = tmp.Write([]byte(bom))
			if err != nil {
				panic(err)
			}
			_, err = tmp.Write(buffer[:n])
			if err != nil {
				panic(err)
			}
			if !endOfFileReached {
				_, err = io.Copy(tmp, inputFile)
				if err != nil {
					panic(err)
				}
			}
			st, err := inputFile.Stat()
			if err != nil {
				panic(err)
			}
			mode := st.Mode()
			err = inputFile.Close()
			if err != nil {
				panic(err)
			}
			outputFile, err := os.OpenFile(fileName, os.O_WRONLY, mode)
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()
			_, err = tmp.Seek(0, 0)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(outputFile, tmp)
			if err != nil {
				panic(err)
			}
			if !quiet {
				fmt.Println("BOM added.")
			}
		} else if shouldStripIfAny && !quiet {
		    if !shouldOutputDetectionResult {
		        //The information is still displayed when it's a no-op.
		        fmt.Println("BOM not found.")
		    }
		    fmt.Println("Nothing is removed.")
		}
	}
}
