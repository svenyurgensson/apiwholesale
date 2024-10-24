package middleware

import (
	"bytes"
	"net/http"
	"time"
	"fmt"
	s "apiwholesale/system"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/mutil"
	 m "github.com/zenazn/goji/web/middleware"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white.
//
// Logger prints a request ID if one is provided.
//
// Logger has been designed explicitly to be Good Enough for use in small
// applications and for people just getting started with Goji. It is expected
// that applications will eventually outgrow this middleware and replace it with
// a custom request logger, such as one that produces machine-parseable output,
// outputs logs to a different service (e.g., syslog), or formats lines like
// those printed elsewhere in the application.
func Logger(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqID := m.GetReqID(*c)

		printStart(reqID, r)

		lw := mutil.WrapWriter(w)

		s.RequestsTotal += 1

		t1 := time.Now()
		h.ServeHTTP(lw, r)

		if lw.Status() == 0 {
			lw.WriteHeader(http.StatusOK)
		}
		t2 := time.Now()

		printEnd(reqID, lw, t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func printStart(reqID string, r *http.Request) {
	var buf bytes.Buffer

	if reqID != "" {
		fmt.Fprintf(&buf, "[%s] ", reqID)
	}
	buf.WriteString("Started ")
	fmt.Fprintf(&buf, "%s", r.Method)
	fmt.Fprintf(&buf, "%q", r.URL.String())
	buf.WriteString("from ")
	buf.WriteString(r.RemoteAddr)

	s.Log.Info(buf.String())
}

func printEnd(reqID string, w mutil.WriterProxy, dt time.Duration) {
	var buf bytes.Buffer

	if reqID != "" {
		fmt.Fprintf(&buf, "[%s] ", reqID)
	}
	buf.WriteString("Returning ")
	status := w.Status()
	fmt.Fprintf(&buf, "%03d", status)
	buf.WriteString(" in ")
	fmt.Fprintf(&buf, "%s", dt)

	if status > 399 {
		s.RequestsFailed += 1
	}

	s.Log.Info(buf.String())
}
