package digits

import "testing"

func TestInit(t *testing.T) {
    var d Digits
    d.Init(1234)
    for di := 1; di <= 4; di++ {
        if d.digits[di - 1] != float64(di) {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di - 1, di, d.digits[di - 1])
        }
    }
}

func TestInitFixed(t *testing.T) {
    var d Digits
    // Normal case first.
    d.InitFixed(4, 1234)
    if len(d.digits) != 4 {
        t.Errorf("Expected len(d) == 4, but got %v", len(d.digits))
    }
    for di := 1; di <= 4; di++ {
        if d.digits[di - 1] != float64(di) {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di - 1, di, d.digits[di - 1])
        }
    }

    // Now test left padding with zeros.
    d.InitFixed(5, 1234)
    if len(d.digits) != 5 {
        t.Errorf("Expected len(d) == 5, but got %v", len(d.digits))
    }
    for di := 0; di <= 4; di++ {
        if d.digits[di] != float64(di) {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di, di, d.digits[di])
        }
    }

    // Now test truncation on the left when the fixed length is too short.
    d.InitFixed(3, 1234)
    if len(d.digits) != 3 {
        t.Errorf("Expected len(d) == 3, but got %v", len(d.digits))
    }
    for di := 2; di <= 4; di++ {
        if d.digits[di - 2] != float64(di) {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di - 2, di, d.digits[di - 2])
        }
    }
}

func TestSetDigit(t *testing.T) {
    var d Digits
    d.Init(4321)
    d.SetDigit(0, 1)
    d.SetDigit(1, 2)
    d.SetDigit(2, 3)
    d.SetDigit(3, 4)

    if len(d.digits) != 4 {
        t.Errorf("Expected len(d) == 4, but got %v", len(d.digits))
    }
    for di := 1; di <= 4; di++ {
        if d.digits[di - 1] != float64(di) {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di - 1, di, d.digits[di - 1])
        }
    }
}

func TestDigit(t *testing.T) {
    var d Digits
    d.Init(1234)
    for di := 1; di <= 4; di++ {
        if d.Digit(di - 1) != float64(di) {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di - 1, di, d.Digit(di - 1))
        }
    }
}

func TestPrint(t *testing.T) {
    // Can't test this one but we'll run it and one can eyeball the results.
    var d Digits
    d.Init(1234)
    d.Print("Test print of 1234")
}

func TestSprint(t *testing.T) {
    var d Digits
    d.Init(1234)
    want := "1234"
    got := d.Sprint()
    if want != got {t.Errorf("Expected \"%v\" but got \"%v\"", want, got)}
}

func TestSprintDigit(t *testing.T) {
    var d Digits
    d.Init(1234)
    want := []string{"1", "2", "3", "4"}
    for di := 0; di <= 3; di++ {
        if d.SprintDigit(di) != want[di] {
            t.Errorf("Expected digit[%v] == %v, but got %v",
                                    di, want[di], d.SprintDigit(di))
        }
    }
}

func TestPartN(t *testing.T) {
    var d Digits
    d.Init(1234)
    want := 12.0
    got := d.PartN(1)
    if want != got {t.Errorf("Expected \"%v\" but got \"%v\"", want, got)}
}

func TestLen(t *testing.T) {
    var d Digits
    d.Init(1234)
    want := 4
    if want != d.Len() {
        t.Errorf("Expected \"%v\" but got \"%v\"", want, d.Len())
    }

    // Pathological case
    d.Init(0)
    want = 1
    if want != d.Len() {
        t.Errorf("Expected \"%v\" but got \"%v\"", want, d.Len())
    }
}

func TestCompare(t *testing.T) {
    var d, d2 Digits
    d.Init(1234)
    d2.Init(1234)
    want := true
    got := d.Compare(d2)
    if want != got {t.Errorf("Expected \"%v\" but got \"%v\"", want, got)}

    d2.Init(4321)
    want = false
    got = d.Compare(d2)
    if want != got {t.Errorf("Expected \"%v\" but got \"%v\"", want, got)}

    d2.Init(123)
    want = false
    got = d.Compare(d2)
    if want != got {t.Errorf("Expected \"%v\" but got \"%v\"", want, got)}
}
