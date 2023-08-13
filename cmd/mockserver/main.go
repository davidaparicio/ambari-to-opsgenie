package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	CONTENT_TYPE = "Content-Type"
	APPLICA_JSON = "application/json"
	JSON_PATH    = "api/examples/"
)

var l *logrus.Logger

type Handler struct {
	StatusCode int
	FilePath   string
}

func jsonHandler(statusCode int, filePath string) Handler {
	return Handler{
		StatusCode: statusCode,
		FilePath:   filePath,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.Debugf("%s called by %s", h.FilePath, r.Host)
	jsonFile, err := os.Open(JSON_PATH + h.FilePath + ".json")

	if err != nil {
		l.WithError(err).Error("cannot open json file")
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			l.WithError(err).Error("cannot close json file")
		}
	}()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		l.WithError(err).Error("cannot json io.ReadAll")
	}
	w.Header().Set(CONTENT_TYPE, APPLICA_JSON)
	w.WriteHeader(h.StatusCode)
	if _, err := w.Write(byteValue); err != nil {
		l.WithError(err).Fatal("ServeHTTP error during w.Write(JSONbytes)")
	}
}

func main() {
	r := http.NewServeMux()
	l = logrus.New()
	l.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05-07:00", FullTimestamp: true})
	l.SetLevel(logrus.DebugLevel)

	r.Handle("/200", jsonHandler(200, "200"))
	r.Handle("/200/critical", jsonHandler(200, "200_critical"))
	r.Handle("/200/warning", jsonHandler(200, "200_warning"))
	r.Handle("/403", jsonHandler(403, "403"))
	r.Handle("/500", jsonHandler(500, "500"))

	srv := &http.Server{
		Addr:              ":1337",
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
		//TLSConfig: tlsConfig,
	}

	l.Debug("Server running on port 1337")
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			l.Debug("Server stopping...")
		} else {
			l.WithError(err).Fatal("Server encontered an unknown error")
		}
	}
}
