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
    _ "embed" // Enables the //go:embed directive for indexHTML.
    "fmt"
    "log"
    "math"
    "net/http"
    "os"
    "strconv"
    "time"
    ld "helloweb/longDiv" // LongDiv produces long division working.
    lm "helloweb/longmult" // LongMult produces long multiplication working.
    "helloweb/prime" // Prime tests for primality.
    quad "helloweb/quadratic" // Quadratic finds the roots.
    series "helloweb/series" // Sum's arithmetic series etc.
)

// index.html is embedded into the binary at build time so that root URL
// requests never read from disk and don't depend on the process working
// directory. Edits to the file take effect only after a rebuild.
// The line below is an ordinary Go comment syntactically, but the build
// toolchain also reads it as a directive: it copies the contents of
// index.html into the variable declared immediately after it.
//go:embed index.html
var indexHTML []byte

// Content type is hard-coded rather than sniffed per request.
const indexContentType = "text/html; charset=utf-8"

// indexContentLength is precomputed once so RootHandler does no per-request
// work to set the header. Setting it explicitly also stops net/http from
// falling back to chunked transfer encoding for this response.
var indexContentLength = strconv.Itoa(len(indexHTML))

// buildETag identifies the source revision index.html was built from. It is
// populated at build time via a linker flag, e.g.:
//   -ldflags "-X main.buildETag=$(git rev-parse --short HEAD)"
// so the value is fixed per binary and never recomputed at runtime. It is
// deliberately short (a git short hash, not a full content hash) so
// If-None-Match comparison stays cheap. It is empty for builds that don't
// pass the flag (e.g. `go run`, `go test`). See GitHub issue #55.
var buildETag string

// maxETagLength is the upper bound from issue #55 criterion 6: the server
// must refuse to start rather than serve an ETag at or beyond this length,
// since a short ETag is the whole point (fast If-None-Match comparison).
const maxETagLength = 128

// activeETag is the ETag RootHandler actually serves, resolved once by
// resolveETag during startup initialisation and held here for the life of
// the process (never recomputed per-request).
var activeETag string

// resolveETag decides the ETag value for this run of the server: buildETag
// when the linker flag supplied one, otherwise (issue #55 criterion 7) the
// number of seconds since the Unix epoch at the moment given by now. It
// also enforces criterion 6's length limit, returning an error that main
// should treat as fatal -- the server must not start while serving an
// ETag that defeats the point of being short.
func resolveETag(buildETag string, now time.Time) (string, error) {
    etag := buildETag
    if etag == "" {
        etag = strconv.FormatInt(now.Unix(), 10)
    }
    if len(etag) > maxETagLength {
        return "", fmt.Errorf("ETag %v exceeds the maximum length allowed of %v characters", etag, maxETagLength)
    }
    return etag, nil
}

func main() {

    if len(indexHTML) == 0 {
        log.Printf("FAIL: embedded index.html is empty\n")
        os.Exit(1)
    }

    etag, err := resolveETag(buildETag, time.Now())
    if err != nil {
        log.Fatal(err)
    }
    activeETag = etag

    // use PORT environment variable, or default to 8080
    port := "8080"
    if fromEnv := os.Getenv("PORT"); fromEnv != "" {
        port = fromEnv
    }

    log.Println("Webserver is starting ...")

    // register hello function to handle all requests
    server := http.NewServeMux()

    server.HandleFunc("/{$}", RootHandler)

    // Explicit catch-all for unmatched paths (e.g. probes from cyber attackers).
    // Keeping this as our own handler, rather than relying on ServeMux's
    // default 404, gives us a place to add logging/metrics on these later.
    server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.NotFound(w, r)
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
    err = http.ListenAndServe(":"+port, server)
    log.Fatal(err)
}

// RootHandler serves the embedded index.html for exact GET requests to "/"
// with no query parameters. Anything else (other methods, or GET with query
// parameters) is rejected immediately, since these are the kinds of
// requests attackers probe with -- this rejection happens unconditionally,
// before ETag negotiation, so a crafted If-None-Match can't bypass it
// (issue #55 criterion 5).
//
// A client whose cached If-None-Match matches activeETag gets a 304 with
// no body; otherwise the response is written straight from memory with
// precomputed headers: no disk I/O, no content sniffing, no range
// handling. This is the leanest path we can offer on what is a common
// attack-probe endpoint.
func RootHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodGet {
        http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
        return
    }

    if len(r.URL.Query()) > 0 {
        http.NotFound(w, r)
        return
    }

    h := w.Header()
    // Set (not the map-literal shortcut used below for the other headers)
    // because net/http canonicalizes "ETag" to "Etag", and Get looks it up
    // under that canonical form. HTTP field names are case-insensitive
    // (RFC 9110 §5.1) -- "ETag" here is written to match the spec's own
    // spelling for readability, not because casing has any wire effect.
    quotedETag := fmt.Sprintf("%q", activeETag)
    h.Set("ETag", quotedETag)

    if r.Header.Get("If-None-Match") == quotedETag {
        w.WriteHeader(http.StatusNotModified)
        return
    }

    h["Content-Type"] = []string{indexContentType}
    h["Content-Length"] = []string{indexContentLength}
    w.Write(indexHTML)
}

// hello responds to the request with a plain-text "Hello, world" message.
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Serving request: %s", r.URL.Path)
    host, _ := os.Hostname()
    fmt.Fprintf(w, "Hello, web!\n")
    fmt.Fprintf(w, "Version: 3.5.1\n")
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

    fmt.Fprintf(w, "%v", prime.CalculateHtml(numberStr))
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

    result := lm.CalculateHtml(number1, number2)

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

    result := ld.CalculateHtml(denom, num)

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
       fmt.Fprintf(w, "%v\n\n:-)", htmlStr)
    default:
        fmt.Fprintf(w, "Rest call /%v/ is not implemented :-(", output)
        fmt.Fprintf(w, " / Found in URL %v", urlStr)
    }
}

