package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

// Provide a test function so this package won't cause an error when
// go test ./... is run from the build root directory.
func TestPrime(t *testing.T) {
    if TestIsPrime() != true {
        t.Errorf("Expected true but got /%v/", false)
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
