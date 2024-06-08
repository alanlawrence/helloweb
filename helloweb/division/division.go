package main

import (
    "os"
    "flag"
    d "helloweb/digits"
    ld "helloweb/longDiv"
)

// This is a small test program which creates a standalone binary.
// It does not itself form part of the webserver any more than
// a *_test.go file does.

// Define the cmd line arguments for this program.
var argDenom = flag.Int("d", 12,
                "d=<denominator> e.g. d=12, must be a +ve integer")
var argNum = flag.Int("n", 143,
                "n=<numerator> e.g. n=143, must be a +ve integer")
var argDebug = flag.Bool("debug", false,
                "debug=<false|true>, turn on/off debug output")

// Globals
var DEBUG = false

func main() {

    flag.Parse()
    if !flag.Parsed() {
        flag.PrintDefaults()
        os.Exit(1)
    }
    DEBUG = *argDebug
    d.DEBUG = *argDebug
    ld.DEBUG = *argDebug
    if !ld.SelfTest() { os.Exit(1)}
    digitsD, quotient, interim, rem := ld.LongDiv(*argDenom, *argNum)
    ld.PrintWorking(digitsD, quotient, interim, rem)
}

