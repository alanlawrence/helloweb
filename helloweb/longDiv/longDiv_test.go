package longDiv

import "testing"

func TestLongDiv(t *testing.T) {

    if !SelfTest() {
        t.Errorf("LongDiv self test failed")
    }

    // Test PrintWorking. Note the stdout is not printed
    // unless there is a test failure, i.e. t.Errorf(...)
    PrintWorking(LongDiv(12, 143))
}

func TestGenerateHtml(t *testing.T) {

    // Test PrintWorking. Note the stdout is not printed
    // unless there is a test failure, i.e. t.Errorf(...)
    PrintWorking(LongDiv( 7,  15))

    htmlStrExp := "   2 remainder=1<br> .--<br>7|15<br>  14<br>  --<br>   1<br>"
    htmlStr := GenerateHtml(LongDiv(7, 15))
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }
}

