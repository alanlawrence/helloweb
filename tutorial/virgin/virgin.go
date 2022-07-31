package main

import (
    "bufio"
    "fmt"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func FixCorruption(myBytes []byte) bool {
    // Process myBytes looking for e3 and replacing with a3, and
    // b<n> with 3<n>, and fix the commas, and decimal points. 
    // No need to contextualise because all occurrences
    // need conversion regardless of location and context.
    // Just loop over the array with a big switch statement and modify
    // in place.
    // Ref: ASCII: https://commons.wikimedia.org/wiki/File:ASCII-Table-wide.svg

    for i:=0; i<len(myBytes); i++ {
        switch myBytes[i] {
        case 0xe3:
            // English £ corrupts so using American pound symbol.
            myBytes[i] = 0x23 // # symbol.
        case 0xb0:
            myBytes[i] = 0x30
        case 0xb1:
            myBytes[i] = 0x31
        case 0xb2:
            myBytes[i] = 0x32
        case 0xb3:
            myBytes[i] = 0x33
        case 0xb4:
            myBytes[i] = 0x34
        case 0xb5:
            myBytes[i] = 0x35
        case 0xb6:
            myBytes[i] = 0x36
        case 0xb7:
            myBytes[i] = 0x37
        case 0xb8:
            myBytes[i] = 0x38
        case 0xb9:
            myBytes[i] = 0x39
        case 0xac:
            myBytes[i] = 0x2c // comma
        case 0xae:
            myBytes[i] = 0x2e // decimal point
        default:
            // Do nothing
        }
    }
    return true
}

func main() {

    myF, err := os.Open("/mnt/c/users/alanj/Downloads/" +
                        "Virgin-transactions-20171004.txt")
    check(err)
    defer myF.Close()
    // Set up reader.
    myReader := bufio.NewReader(myF)

    //Output file for conversion.
    outF := "converted.txt"
    myFw, err := os.Create(outF)
    check(err)
    defer myFw.Close()
    // Set up writer.
    myWriter := bufio.NewWriter(myFw)

    // Run a loop reading an array of bytes from the file
    // and write those bytes to another file.

    // Number of bytes written in the Write(bytes) call.
    var numBytesW int
    // Total number of bytes written
    numBytesWTotal := 0
    // Array of bytes to read into.
    bufSize := 64
    myBytes := make([]byte, bufSize)
    // Number of bytes read in the Read(num bytes) call.
    numBytesR := 1 // Artificial value to enter the loop.

    for numBytesR > 0 {
        // Read from the input file. I guess the number read depends on the
        // size of the Reader's buffer.
        numBytesR, err = myReader.Read(myBytes)
        if numBytesR > 0 {
            check(err)  // We will get a legitimate error at EOF when
                        // read bytes = 0.

            FixCorruption(myBytes)

            // Write myBytes to the output file.
            // Pass a slice of the correct size.
            numBytesW, err = myWriter.Write(myBytes[:numBytesR])


            // Report progress.
            numBytesWTotal += numBytesW
            check(err)
            fmt.Printf("wrote %5d bytes\r\r\r\r\r\r\r\r\r\r\r", numBytesWTotal)
            myWriter.Flush()
        }
    }
    // Report final progress.
    fmt.Printf("wrote %5d bytes\n", numBytesWTotal)

}

