package series

import (
    "testing"
    "net/url"
    "fmt"
)

// Note: these tests are prone to floating point arithmetic anomalies.
// TODO: use tolerance checks.
func TestSeries(t *testing.T) {

    a := 1.0
    d := 1.0
    n := 10.0

    arSum, err := ARSum(a, d, n)
    if arSum != 55.0 {
        t.Errorf("Expected arSum=55, but got %v with err: /%v/", arSum, err)
    }

    a = 10.0
    d = -1.0
    n = 10.0

    arSum, err = ARSum(a, d, n)
    if arSum != 55.0 {
        t.Errorf("Expected arSum=55, but got %v with err: /%v/", arSum, err)
    }

    a =  0.0
    d = -1.0
    n = 5.0

    arSum, err = ARSum(a, d, n)
    if arSum != -10.0 {
        t.Errorf("Expected arSum=-10, but got %v with err: /%v/", arSum, err)
    }

    a = -9.0
    d =  1.0
    n = 10.0

    arSum, err = ARSum(a, d, n)
    if arSum != -45.0 {
        t.Errorf("Expected arSum=-45, but got %v with err: /%v/", arSum, err)
    }

    a = 0.0
    d = 0.5
    n = 10.0

    arSum, err = ARSum(a, d, n)
    if arSum != 22.5 {
        t.Errorf("Expected arSum=22.5, but got %v with err: /%v/", arSum, err)
    }

    // For this one, we get a floating point anomaly on my local machine
    // and also on The Go Playground, so it is probably a stable anomaly.
    a = 0.5
    d = -0.1
    n = 10.0

    arSum, err = ARSum(a, d, n)
    epsilon := 1e-6  // Tolerance for floating point anomaly.
    if !(arSum > 0.5 - epsilon && arSum < 0.5 + epsilon) {
        t.Errorf("Expected arSum=0.5, but got %v with err: /%v/", arSum, err)
    }

    a = 0.0
    d = 1.0
    n = 0.0

    errExp := "Expected n >= 1 but got 0"
    arSum, err = ARSum(a, d, n)
    if (arSum != 0.0 || err != errExp) {
        t.Errorf("Expected arSum=0, err: /%v/, but got %v with err: /%v/",
                  errExp, arSum, err)
    }
}

func TestGenerateHtml(t *testing.T) {

    htmlStrExp := "For a=1, d=1, n=10, the sum is 55"
    htmlStr := GenerateHtml(1.0, 1.0, 10, 55.0, "")
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    htmlStrExp = "For a=1, d=1, n=10, the sum is 54<br>Oops! blast!"
    htmlStr = GenerateHtml(1.0, 1.0, 10, 54.0, "blast!")
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }
}

func TestARCalc(t *testing.T) {
    // Good case.
    urlStr := "a=1&d=1&n=10"
    htmlStrExp := "For a=1, d=1, n=10, the sum is 55"

    // Initialise values
    v, _ := url.ParseQuery(urlStr)

    // Call parse function to check output.
    htmlStr := ARCalc(urlStr, v)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    // Pathological case: param n is mising
    urlStr = "a=1&d=1"
    htmlStrExp = fmt.Sprintf("Expected 3 params /a, d, n/ but got %v", urlStr)

    // Initialise values
    v, _ = url.ParseQuery(urlStr)

    // Call parse function to check output.
    htmlStr = ARCalc(urlStr, v)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    // Pathological case: param n is not set
    // Add a test when we have implemented the desired behaviour.
}


func TestDebug(t *testing.T) {
    // Set DEBUG to true to cover any debug lines.
    debugBefore := DEBUG
    DEBUG=true
    a := 1.0
    d := 1.0
    n := 10.0

    arSum, err := ARSum(a, d, n)
    if arSum != 55.0 {
        t.Errorf("Expected arSum=55, but got %v with err: /%v/", arSum, err)
    }
    DEBUG = debugBefore
}
