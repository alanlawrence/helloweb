/*
 * This webserver accepts some simple input and returns a response.
 * Entry point is URL/  // Which finds the index.html file with the number form.
 * The form invokes an ajax request to URL/prime?number=13.
 * Or one can access the prime number checker directly by typing 
 *    URL/prime?number=13
 * directly in the browser address bar. It is all the same to this webserver.
 * The /prime handler extracts the number and works out if it is prime and 
 * writes the result back.
 * Other functions are similarly implemented.
 */

package main

import (
    "fmt"
    "log"
    "math"
    "net/http"
    "os"
    "strconv"
    "time"
    ld "helloweb/longDiv" // LongDiv produces long division working.
    quad "helloweb/quadratic" // Quadratic finds the roots.
    series "helloweb/series" // Sum's arithmetic series etc.
)

func main() {

    // Test functions and don't run the webserver if the tests fail.
    fmt.Printf("Running self tests ... ")
    if !TestIsPrime() {
        os.Exit(1)
    }
    // TODO: what about the other self test functions?
    //     : those in packages are tested with "go test <package>"
    fmt.Printf("passed. Starting webserver ...\n")


    // use PORT environment variable, or default to 8080
    port := "8080"
    if fromEnv := os.Getenv("PORT"); fromEnv != "" {
        port = fromEnv
    }

    // register hello function to handle all requests
    server := http.NewServeMux()

    // The default page is served from the local directory.
    server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })

    server.HandleFunc("/hello", hello)
    server.HandleFunc("/prime", PrimeHandler)
    server.HandleFunc("/gcd", GcdHandler)
    server.HandleFunc("/longmult", LongMultHandler)
    server.HandleFunc("/longdiv", LongDivHandler)
    server.HandleFunc("/quadratic", QuadraticHandler)
    server.HandleFunc("/ar-series/{output}", ArSeriesHandler)



    // start the web server on port and accept requests
    fmt.Printf("Server listening on port %s\n", port)
    err := http.ListenAndServe(":"+port, server)
    log.Fatal(err)
}

// hello responds to the request with a plain-text "Hello, world" message.
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Serving request: %s", r.URL.Path)
    host, _ := os.Hostname()
    fmt.Fprintf(w, "Hello, GKE!\n")
    fmt.Fprintf(w, "Version: 3.1.0\n")
    fmt.Fprintf(w, "Hostname: %s\n", host)
    fmt.Fprintf(w, "Private message: Daddy loves you Pops!\n")
    fmt.Fprintf(w, "Time: %v\n", time.Now())
}

func PrimeHandler(w http.ResponseWriter, r *http.Request) {

    numbers, ok := r.URL.Query()["number"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'number' is missing")
        return
    }

    // Query()["number"] will return an array of items, 
    // we only want the single item.
    numberStr := numbers[0]

    number, _ := strconv.Atoi(numberStr)

    result := ""
    if IsPrime(number) {
        result = "is prime!"
    } else {
        result = "is not prime :-("
    }

    fmt.Fprintf(w, "%v %v", number, result)
}

// TODO: Move IsPrime and its test driver out into another package.
func IsPrime(n int) bool {

    isPrime := true
    done := false
    if n <= 3 {
        isPrime = n >= 1
        done = true
    } else if (n % 2 == 0 || n % 3 == 0) {
        isPrime = false
        done = true
    }

    // All the numbers below 25 are either divisible by 2 or 3 (tested above)
    // or are prime.
    if (!done && n < 25) {
        isPrime = true
        done = true
    }

    // Now exploit the property that all primes >=6 are of the form 6k+1 or 6k-1
    // since 2 divides 6k, 6k+2 and 6k+4, and 3 divides 6k+3
    // which leaves 6k+1 and 6k+5 (== 6k'-1, where k'=k+1).

    // So we test all numbers of the form 6k+/-1 such that
    //     6k+/-1     <= sqrt(n)
    // ==> (6k+/-1)^2 <= n

    i := 5
    // This generates the pair 5 and 7 for the first iteration, k = 1
    for (!done && i*i <= n) {

        if n % i == 0 {
            isPrime = false
            done = true
        } else if n % (i+2) == 0 {
            isPrime = false
            done = true
        }
        // Advance to next iteration. Imagine k += 1
        i += 6
    }
    // else must be prime hence return isPrime default of true.

    return isPrime
}

func TestIsPrime() bool {

    pass := 0
    tests := 0
    if IsPrime(1) && IsPrime(2) && IsPrime(3) && IsPrime(5) && IsPrime(7) && IsPrime(11) {
        pass += 6
    } else {
        log.Printf("FAIL: isPrime returned false for a prime\n")
    }
    tests += 6

    if !(IsPrime(4) || IsPrime(6) || IsPrime(8) || IsPrime(9) || IsPrime(10)) {
        pass += 5
    } else {
        log.Printf("FAIL: isPrime returned true for a non-prime\n")
    }
    tests += 5

    // Test some larger numbers to exercise the k+/-1 loop.
    if IsPrime(23) && IsPrime(29) && IsPrime(37) && IsPrime(97) && IsPrime(107) {
        pass += 5
    } else {
        log.Printf("FAIL: isPrime returned false for a prime\n")
    }
    tests += 5

    if !(IsPrime(24) || IsPrime(25) || IsPrime(38) || IsPrime(99) || IsPrime(115)) {
        pass += 5
    } else {
        log.Printf("FAIL: isPrime returned true for a non-prime\n")
    }
    tests += 5

    return pass == tests
}

func GcdHandler(w http.ResponseWriter, r *http.Request) {

    numbers, ok := r.URL.Query()["number1"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'number1' is missing")
        return
    }

    // Query()["number"] will return an array of items, 
    // we only want the single item.
    number1Str := numbers[0]

    number1, _ := strconv.ParseFloat(number1Str, 64)

    numbers, ok = r.URL.Query()["number2"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'number2' is missing")
        return
    }

    // Query()["number"] will return an array of items, 
    // we only want the single item.
    number2Str := numbers[0]

    number2, _ := strconv.ParseFloat(number2Str, 64)

    result := Gcd(number1, number2) 

    if result > 1 {
        fmt.Fprintf(w, "%v :-)", result)
    } else {
        fmt.Fprintf(w, "Nothing other than 1 :-(")
    }
}

func Gcd(ra float64, rb float64) float64 {

    rt := 0.0
    for rb != 0 {
        rt = rb
        rb = math.Mod(ra, rb)
        ra = rt
    }
    return ra
}

func LongMultHandler(w http.ResponseWriter, r *http.Request) {

    numbers, ok := r.URL.Query()["number1"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'number1' is missing")
        return
    }

    // Query()["number"] will return an array of items, 
    // we only want the single item.
    number1Str := numbers[0]

    number1, _ := strconv.ParseFloat(number1Str, 64)

    numbers, ok = r.URL.Query()["number2"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'number2' is missing")
        return
    }

    // Query()["number"] will return an array of items, 
    // we only want the single item.
    number2Str := numbers[0]

    number2, _ := strconv.ParseFloat(number2Str, 64)

    result := GenerateHtml(LongMult(int(number1), int(number2)))

    fmt.Fprintf(w, "%v\n\n:-)", result)
}

func LongDivHandler(w http.ResponseWriter, r *http.Request) {

    numbers, ok := r.URL.Query()["denom"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'denom' is missing")
        return
    }

    // Query()["denom"] will return an array of items, 
    // we only want the single item.
    number1Str := numbers[0]

    denom, _ := strconv.ParseFloat(number1Str, 64)

    numbers, ok = r.URL.Query()["num"]

    if !ok || len(numbers[0]) < 1 {
        log.Println("Url Param 'num' is missing")
        return
    }

    // Query()["num"] will return an array of items, 
    // we only want the single item.
    number2Str := numbers[0]

    num, _ := strconv.ParseFloat(number2Str, 64)

    result := ld.GenerateHtml(ld.LongDiv(int(denom), int(num)))

    fmt.Fprintf(w, "%v\n\n:-)", result)
}

func QuadraticHandler(w http.ResponseWriter, r *http.Request) {

    coeffs, ok := r.URL.Query()["a"]

    if !ok || len(coeffs[0]) < 1 {
        log.Println("Url Param 'a' is missing")
        return
    }

    // Query()["a"] will return an array of items, 
    // we only want the single item.
    coeffAStr := coeffs[0]

    coeffA, _ := strconv.ParseFloat(coeffAStr, 64)

    coeffs, ok = r.URL.Query()["b"]

    if !ok || len(coeffs[0]) < 1 {
        log.Println("Url Param 'b' is missing")
        return
    }

    // Query()["b"] will return an array of items, 
    // we only want the single item.
    coeffBStr := coeffs[0]

    coeffB, _ := strconv.ParseFloat(coeffBStr, 64)

    coeffs, ok = r.URL.Query()["c"]

    if !ok || len(coeffs[0]) < 1 {
        log.Println("Url Param 'c' is missing")
        return
    }

    // Query()["c"] will return an array of items, 
    // we only want the single item.
    coeffCStr := coeffs[0]

    coeffC, _ := strconv.ParseFloat(coeffCStr, 64)

    root1, root2 := quad.Quadratic(coeffA, coeffB, coeffC)
    result := quad.GenerateHtml(coeffA, coeffB, coeffC, root1, root2)

    fmt.Fprintf(w, "%v\n\n:-)", result)
}

func ArSeriesHandler(w http.ResponseWriter, r *http.Request) {

    urlStr := r.URL.String()
    // Retrieve value of {output}
    output := r.PathValue("output")
    switch output {
    case "sum":
       htmlStr := series.ARCalc(urlStr, r.URL.Query())
       fmt.Fprintf(w, htmlStr)
    default:
        fmt.Fprintf(w, "Rest call /%v/ is not implemented :-(", output)
        fmt.Fprintf(w, " / Found in URL %v", urlStr)
    }
}

// Algorithm: https://en.wikipedia.org/wiki/Multiplication_algorithm#Long_multiplication
func LongMult(bnum int, anum int) ([]float64, []float64,
                                   [][]float64, [][]float64) {

    base := 10.0
    //turn ints into a string of int digits
    lenA := 1
    if anum > 0 {
        lenA = int(math.Trunc(math.Log10(float64(anum)) + 1.0))
    }
    lenB := 1
    if bnum > 0 {
        lenB = int(math.Trunc(math.Log10(float64(bnum)) + 1.0))
    }
    digitsA := make([]float64,lenA)
    for di := 0; di < lenA; di++ {
        digitsA[di] = math.Mod(
            math.Trunc(float64(anum)/math.Pow10(lenA-1-di)), base)
    }
    digitsB := make([]float64,lenB)
    for di := 0; di < lenB; di++ {
        digitsB[di] = math.Mod(
            math.Trunc(float64(bnum)/math.Pow10(lenB-1-di)), base)
    }

    //product := array [lenB + 1][lenA + lenB]
    //    A row of interim results for each digit in multiplier bnum
    //    and 1 extra row to hold the Total.
    //    The interim result rows are not necessary to compute the Total
    //    but they are required in long multiplication working.
    product := make([][]float64, lenB + 1)
    for ri := range product {
        product[ri] = make([]float64, lenA + lenB)
    }
    //An array to store the mulitplication carries. 
    //Not necessary for the calculation
    //but required in long multiplication working.
    carrys := make([][]float64, lenB)
    for ri := range carrys {
        carrys[ri] = make([]float64, lenA + lenB)
    }

    carryTotal := 0.0
    carryRow := 0.0
    pi := 0
    // We descend the array from left to right so we read
    // the digits in human readable order.
    for bi := lenB-1; bi >= 0; bi-- {
        carryTotal = 0
        carryRow = 0
        for ai := lenA-1; ai >= 0; ai-- {
            pi = ai + bi + 1
            // Per row interim results (row per digit of mulitplier anum)
            product[bi][pi] = carryRow + digitsA[ai] * digitsB[bi]
            carryRow = math.Trunc(product[bi][pi] / base)
            carrys[bi][pi] = carryRow
            product[bi][pi] = math.Mod(product[bi][pi], base)

            // Total. Note the += to sum the interim results as we go.
            product[lenB][pi] += carryTotal + digitsA[ai] * digitsB[bi]
            carryTotal = math.Trunc(product[lenB][pi] / base)
            product[lenB][pi] = math.Mod(product[lenB][pi], base)
        }
        product[bi][bi] = carryRow
        product[lenB][bi] = carryTotal
    }
    return digitsA, digitsB, product, carrys
}

func GenerateHtml(digitsA []float64, digitsB []float64,
                  product [][]float64, carrys [][]float64) string {

    workingStr := fmt.Sprintf("")
    lenS := len(digitsA) - len(digitsB)
    for si := 0; si < -1 * lenS; si++ {
        workingStr += fmt.Sprintf("  ")
    }
    workingStr += fmt.Sprintf("            %v<br>", digitsA)
    for si := 0; si < lenS; si++ {
        workingStr += fmt.Sprintf("  ")
    }
    workingStr += fmt.Sprintf("          x %v<br><br>", digitsB)
    lenP := len(product)
    printZeroes := false
    for ri := lenP - 2 ; ri >= 0; ri-- {
        workingStr += fmt.Sprintf("          ")
        printZeroes = false
        for _, digitP := range product[ri] {
            if digitP == 0 && !printZeroes {
                workingStr += fmt.Sprintf("  ")
            } else {
                workingStr += fmt.Sprintf("%v ", digitP)
                printZeroes = true
            }
        }
        workingStr += fmt.Sprintf("<br>");
        workingStr += fmt.Sprintf("carry-> ")
        for _, digitC := range carrys[ri] {
            if digitC == 0 {
                workingStr += fmt.Sprintf("  ")
            } else {
                workingStr += fmt.Sprintf("%v ", digitC)
            }
        }
        workingStr += fmt.Sprintf("<br>");
    }
    workingStr += fmt.Sprintf("Total     ")
    printZeroes = false
    for _, digitT := range product[lenP-1] {
        if digitT == 0 && !printZeroes {
            workingStr += fmt.Sprintf("  ")
        } else {
            workingStr += fmt.Sprintf("%v ", digitT)
            printZeroes = true
        }
    }
    workingStr += fmt.Sprintf("<br>");

    return workingStr
}

