package main

import (
    "net/http"
    "net/http/httptest"
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
