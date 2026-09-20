package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
)

// TestPrimeHandler covers the /prime endpoint's query-param wiring;
// primality logic itself is covered by helloweb/prime's tests.
func TestPrimeHandler(t *testing.T) {
    tests := []struct {
        name   string
        target string
        want   string
    }{
        {"valid prime", "/prime?number=13", "13 is prime!"},
        {"valid composite", "/prime?number=4", "4 is not prime :-("},
        {"zero is invalid", "/prime?number=0", "Only positive integers are valid inputs"},
        {"negative is invalid", "/prime?number=-5", "Only positive integers are valid inputs"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tc.target, nil)
            rec := httptest.NewRecorder()

            PrimeHandler(rec, req)

            body, err := io.ReadAll(rec.Result().Body)
            if err != nil {
                t.Fatalf("reading response body: %v", err)
            }
            if got := string(body); got != tc.want {
                t.Errorf("PrimeHandler(%s) body = %q, want %q", tc.target, got, tc.want)
            }
        })
    }
}

// TestHyphenHandler covers the /hyphen endpoint's query-param wiring;
// hyphenation logic itself is covered by helloweb/hyphenate's tests.
func TestHyphenHandler(t *testing.T) {
    tests := []struct {
        name   string
        target string
        want   string
    }{
        {"simple string", "/hyphen?text=My%20file", "My-file\n\n:-)"},
        {"missing param", "/hyphen", ""},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tc.target, nil)
            rec := httptest.NewRecorder()

            HyphenHandler(rec, req)

            body, err := io.ReadAll(rec.Result().Body)
            if err != nil {
                t.Fatalf("reading response body: %v", err)
            }
            if got := string(body); got != tc.want {
                t.Errorf("HyphenHandler(%s) body = %q, want %q", tc.target, got, tc.want)
            }
        })
    }
}

// TestHyphenHandlerRejectsOverMaxLength covers the issue #56 criterion 4
// error path through the handler's query-param wiring.
func TestHyphenHandlerRejectsOverMaxLength(t *testing.T) {
    longText := ""
    for i := 0; i < 257; i++ {
        longText += "a"
    }
    req := httptest.NewRequest(http.MethodGet, "/hyphen?text="+longText, nil)
    rec := httptest.NewRecorder()

    HyphenHandler(rec, req)

    body, err := io.ReadAll(rec.Result().Body)
    if err != nil {
        t.Fatalf("reading response body: %v", err)
    }
    want := "Strings submitted for hyphenation must be less than " +
        "or equal to 256 characters\n\n:-("
    if got := string(body); got != want {
        t.Errorf("HyphenHandler body = %q, want %q", got, want)
    }
}

func TestRootHandler(t *testing.T) {
    tests := []struct {
        name       string
        method     string
        target     string
        wantStatus int
    }{
        {"GET with no query parameters serves index.html", http.MethodGet, "/", http.StatusOK},
        {"non-GET method rejected", http.MethodPost, "/", http.StatusMethodNotAllowed},
        {"GET with query parameters rejected", http.MethodGet, "/?number=13", http.StatusNotFound},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(tc.method, tc.target, nil)
            rec := httptest.NewRecorder()

            RootHandler(rec, req)

            if rec.Code != tc.wantStatus {
                t.Errorf("RootHandler(%s %s) status = %v, want %v", tc.method, tc.target, rec.Code, tc.wantStatus)
            }
        })
    }
}

// TestRootHandlerServesEmbeddedContent checks that a plain GET / returns the
// embedded index.html verbatim, with the precomputed headers.
func TestRootHandlerServesEmbeddedContent(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %v, want %v", rec.Code, http.StatusOK)
    }
    if !bytes.Equal(rec.Body.Bytes(), indexHTML) {
        t.Errorf("body does not match embedded index.html")
    }
    if got := rec.Header().Get("Content-Type"); got != indexContentType {
        t.Errorf("Content-Type = %q, want %q", got, indexContentType)
    }
    if got := rec.Header().Get("Content-Length"); got != indexContentLength {
        t.Errorf("Content-Length = %q, want %q", got, indexContentLength)
    }
}

// --- Build-time ETag / conditional GET (GitHub issue #55) ---
//
// activeETag is the resolved, effective ETag RootHandler serves (see its
// doc comment in webserver.go). These request-level tests set it directly
// to simulate a specific run of the server; resolveETag itself -- which
// derives activeETag from buildETag at startup -- gets its own unit tests
// further down, since it's a pure function.

// withActiveETag sets activeETag for the duration of a test and restores
// the previous value on cleanup, so test order doesn't matter.
func withActiveETag(t *testing.T, value string) {
    t.Helper()
    orig := activeETag
    activeETag = value
    t.Cleanup(func() { activeETag = orig })
}

// TestRootHandlerNoIfNoneMatch covers acceptance criterion 2 of issue #55:
// the full page is served when the client sends no If-None-Match at all,
// and the response still carries the server's ETag for the client to cache
// for next time.
func TestRootHandlerNoIfNoneMatch(t *testing.T) {
    withActiveETag(t, "abc1234")

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %v, want %v", rec.Code, http.StatusOK)
    }
    if !bytes.Equal(rec.Body.Bytes(), indexHTML) {
        t.Errorf("body does not match embedded index.html")
    }
    if got, want := rec.Header().Get("ETag"), `"abc1234"`; got != want {
        t.Errorf("ETag = %q, want %q", got, want)
    }
}

// TestRootHandlerIfNoneMatchMismatch covers acceptance criterion 3: the
// full page is still served when If-None-Match is present but doesn't
// match the server's current ETag.
func TestRootHandlerIfNoneMatchMismatch(t *testing.T) {
    withActiveETag(t, "abc1234")

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("If-None-Match", `"def5678"`)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %v, want %v", rec.Code, http.StatusOK)
    }
    if !bytes.Equal(rec.Body.Bytes(), indexHTML) {
        t.Errorf("body does not match embedded index.html")
    }
}

// TestRootHandlerIfNoneMatchStaleAcrossRebuild is the restart/redeploy edge
// case: a browser cached an ETag from a previous build (previous git
// commit). The server is then rebuilt and redeployed with a new
// buildETag -- there is no per-request or per-connection state carried
// across that, so a client's pre-redeploy ETag must be treated the same as
// any other mismatch: full page served, never a wrongly-cached 304.
func TestRootHandlerIfNoneMatchStaleAcrossRebuild(t *testing.T) {
    withActiveETag(t, "newcommit")

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("If-None-Match", `"oldcommit"`) // cached from before the redeploy
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %v, want %v", rec.Code, http.StatusOK)
    }
    if !bytes.Equal(rec.Body.Bytes(), indexHTML) {
        t.Errorf("body does not match embedded index.html")
    }
    if got, want := rec.Header().Get("ETag"), `"newcommit"`; got != want {
        t.Errorf("ETag = %q, want %q", got, want)
    }
}

// TestRootHandlerIfNoneMatchMatch covers acceptance criterion 4: a matching
// If-None-Match short-circuits to 304 with no body.
func TestRootHandlerIfNoneMatchMatch(t *testing.T) {
    withActiveETag(t, "abc1234")

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("If-None-Match", `"abc1234"`)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusNotModified {
        t.Fatalf("status = %v, want %v", rec.Code, http.StatusNotModified)
    }
    if rec.Body.Len() != 0 {
        t.Errorf("body length = %v, want 0", rec.Body.Len())
    }
    if got, want := rec.Header().Get("ETag"), `"abc1234"`; got != want {
        t.Errorf("ETag = %q, want %q", got, want)
    }
}

// TestRootHandlerETagIsShort covers the issue's performance requirement
// that the served ETag stay short (<128 characters, quotes included) so
// If-None-Match comparison is cheap.
func TestRootHandlerETagIsShort(t *testing.T) {
    withActiveETag(t, "abc1234")

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if got := rec.Header().Get("ETag"); len(got) >= 128 {
        t.Errorf("ETag length = %v, want < 128 (value: %q)", len(got), got)
    }
}

// TestRootHandlerQueryParamsRejectedRegardlessOfETagMatch covers acceptance
// criterion 5 as clarified: a root URL request with query params ("bad
// params") is still rejected with 404, even when the request also carries
// an If-None-Match value that matches the server's ETag exactly. The
// existing query-param rejection (see TestRootHandler) must take priority
// over ETag negotiation, not the other way round.
func TestRootHandlerQueryParamsRejectedRegardlessOfETagMatch(t *testing.T) {
    withActiveETag(t, "abc1234")

    req := httptest.NewRequest(http.MethodGet, "/?number=13", nil)
    req.Header.Set("If-None-Match", `"abc1234"`)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusNotFound {
        t.Errorf("status = %v, want %v", rec.Code, http.StatusNotFound)
    }
}

// --- resolveETag (GitHub issue #55 criteria 6 and 7) ---
//
// resolveETag is a pure function so startup validation/defaulting can be
// unit-tested without touching global state or exercising main(), which
// exits the process on failure and so isn't practical to test directly.

// TestResolveETagUsesBuildETagWhenPresent is the common case: a non-empty
// buildETag (as set by the linker flag) is used as-is.
func TestResolveETagUsesBuildETagWhenPresent(t *testing.T) {
    got, err := resolveETag("abc1234", time.Now())
    if err != nil {
        t.Fatalf("resolveETag returned unexpected error: %v", err)
    }
    if got != "abc1234" {
        t.Errorf("resolveETag = %q, want %q", got, "abc1234")
    }
}

// TestResolveETagDefaultsToEpochSecondsWhenBuildETagEmpty covers
// acceptance criterion 7: when no build-time ETag was supplied, the
// server defaults to the number of seconds since the Unix epoch, fixed at
// startup.
func TestResolveETagDefaultsToEpochSecondsWhenBuildETagEmpty(t *testing.T) {
    now := time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
    want := fmt.Sprintf("%d", now.Unix())

    got, err := resolveETag("", now)
    if err != nil {
        t.Fatalf("resolveETag returned unexpected error: %v", err)
    }
    if got != want {
        t.Errorf("resolveETag = %q, want %q", got, want)
    }
}

// TestResolveETagRejectsOversizedETag covers acceptance criterion 6: a
// build-supplied ETag of more than 128 characters must be rejected with
// an error carrying the exact wording the issue specifies, so the caller
// (main) can log it and refuse to start.
func TestResolveETagRejectsOversizedETag(t *testing.T) {
    oversized := make([]byte, maxETagLength+1)
    for i := range oversized {
        oversized[i] = 'a'
    }
    tooLong := string(oversized)

    _, err := resolveETag(tooLong, time.Now())
    if err == nil {
        t.Fatal("resolveETag returned no error for an oversized ETag, want one")
    }

    want := fmt.Sprintf("ETag %v exceeds the maximum length allowed of %v characters", tooLong, maxETagLength)
    if got := err.Error(); got != want {
        t.Errorf("resolveETag error = %q, want %q", got, want)
    }
}

// TestResolveETagAllowsExactly128Chars pins the boundary: criterion 6 only
// rejects an ETag *exceeding* 128 characters, so exactly 128 must be
// accepted.
func TestResolveETagAllowsExactly128Chars(t *testing.T) {
    exact := make([]byte, maxETagLength)
    for i := range exact {
        exact[i] = 'a'
    }
    boundary := string(exact)

    got, err := resolveETag(boundary, time.Now())
    if err != nil {
        t.Fatalf("resolveETag returned unexpected error for a %d-char ETag: %v", maxETagLength, err)
    }
    if got != boundary {
        t.Errorf("resolveETag = %q, want %q", got, boundary)
    }
}

// TestRootHandlerIsCwdIndependent proves the handler no longer depends on the
// process working directory containing index.html.
func TestRootHandlerIsCwdIndependent(t *testing.T) {
    orig, err := os.Getwd()
    if err != nil {
        t.Fatalf("Getwd: %v", err)
    }
    t.Cleanup(func() { os.Chdir(orig) })
    if err := os.Chdir(t.TempDir()); err != nil {
        t.Fatalf("Chdir: %v", err)
    }

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()

    RootHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %v, want %v", rec.Code, http.StatusOK)
    }
    if !bytes.Equal(rec.Body.Bytes(), indexHTML) {
        t.Errorf("body does not match embedded index.html")
    }
}
