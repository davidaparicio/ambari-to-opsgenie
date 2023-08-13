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
	CONTENT_TYPE      = "Content-Type"
	APPLICA_JSON      = "application/json"
	JSON_PATH         = "api/examples/"
	READTIMEOUT       = 1
	WRITETIMEOUT      = 1
	IDLETIMEOUT       = 30
	READHEADERTIMEOUT = 2
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

	r.Handle("/200", jsonHandler(http.StatusOK, "200"))
	r.Handle("/200/critical", jsonHandler(http.StatusOK, "200_critical"))
	r.Handle("/200/warning", jsonHandler(http.StatusOK, "200_warning"))
	r.Handle("/403", jsonHandler(http.StatusForbidden, "403"))
	r.Handle("/500", jsonHandler(http.StatusInternalServerError, "500"))

	srv := &http.Server{
		Addr:              ":1337",
		ReadTimeout:       READTIMEOUT * time.Second,
		WriteTimeout:      WRITETIMEOUT * time.Second,
		IdleTimeout:       IDLETIMEOUT * time.Second,
		ReadHeaderTimeout: READHEADERTIMEOUT * time.Second,
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
