package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type testCase struct {
	name             string
	method           string
	url              string //not useful for TestJsonHandler
	statusExpected   int
	containsExpected string
}

func init() {
	// Initialize the logrus logger for the server
	l = logrus.New()
	l.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05-07:00", FullTimestamp: true})
	l.SetLevel(logrus.ErrorLevel) //l.SetLevel(logrus.DebugLevel)

	// Move to the root project directory to reach ./api/examples/*
	// https://brandur.org/fragments/test-go-project-root
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		l.Panicf("could not change folder: %s", err)
	}
}

func TestJsonHandler(t *testing.T) {
	tt := []testCase{
		{"200 OK", http.MethodGet, "200", http.StatusOK, "OK"},
		{"200 Critical", http.MethodGet, "200_critical", http.StatusOK, "CRITICAL"},
		{"200 Warning", http.MethodGet, "200_warning", http.StatusOK, "WARNING"},
		{"403", http.MethodGet, "403", http.StatusForbidden, "Missing authentication token"},
		{"405 HEAD", http.MethodHead, "", http.StatusMethodNotAllowed, ""},
		{"405 POST", http.MethodPost, "", http.StatusMethodNotAllowed, ""},
		{"405 PUT", http.MethodPut, "", http.StatusMethodNotAllowed, ""},
		{"405 PATCH", http.MethodPatch, "", http.StatusMethodNotAllowed, ""},
		{"405 DELETE", http.MethodDelete, "", http.StatusMethodNotAllowed, ""},
		{"405 CONNECT", http.MethodConnect, "", http.StatusMethodNotAllowed, ""},
		{"405 OPTIONS", http.MethodOptions, "", http.StatusMethodNotAllowed, ""},
		{"405 TRACE", http.MethodTrace, "", http.StatusMethodNotAllowed, ""},
		{"405 Rickroolling", http.MethodPost, "Rick_Astley/Never_Gonna_Give_You_Up", http.StatusMethodNotAllowed, ""},
		{"500", http.MethodGet, "500", http.StatusInternalServerError, "org.apache.ambari.server.controller.spi.SystemException"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "not used/jsonHandler called directly, testRouting will test URLs", nil)
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			rec := httptest.NewRecorder()

			jsonHandler(tc.statusExpected, tc.url).ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			tc.verify(t, res)
		})
	}
}

func TestRouting(t *testing.T) {
	tt := []testCase{
		{"200 OK", http.MethodGet, "200", http.StatusOK, "OK"},
		{"200 Critical", http.MethodGet, "200/critical", http.StatusOK, "CRITICAL"},
		{"200 Warning", http.MethodGet, "200/warning", http.StatusOK, "WARNING"},
		{"403", http.MethodGet, "403", http.StatusForbidden, "Missing authentication token"},
		{"404 Rickroolling", http.MethodPost, "Rick_Astley/Never_Gonna_Give_You_Up", http.StatusNotFound, ""},
		{"500", http.MethodGet, "500", http.StatusInternalServerError, "org.apache.ambari.server.controller.spi.SystemException"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(handler())
			defer srv.Close()

			res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, tc.url))
			if err != nil {
				t.Fatalf("could not GET request: %v", err)
			}
			//t.Logf("%s/%s", srv.URL, tc.url)
			tc.verify(t, res)
		})
	}
}

func (tc testCase) verify(t *testing.T, res *http.Response) {
	if res.StatusCode != tc.statusExpected {
		t.Errorf("expected status %v; got %v", tc.statusExpected, res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	if !strings.Contains(string(b), tc.containsExpected) {
		t.Fatalf("doesn't contain this sentence: %s", tc.containsExpected)
	}
}
