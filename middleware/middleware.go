package middleware

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/getsentry/raven-go"
	"net/http"
	"runtime"
	"time"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger     *logrus.Logger
	Raven      *raven.Client
	PrintStack bool
	StackAll   bool
	StackSize  int
}

type Logger struct {
	*logrus.Logger
}

func trace() *raven.Stacktrace {
	return raven.NewStacktrace(0, 2, nil)
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); r != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			if rec.Raven != nil {
				er := errors.New(fmt.Sprintf("%v", err))

				packet := raven.NewPacket(er.Error(), raven.NewException(er, trace()), raven.NewHttp(r))
				eventID, _ := rec.Raven.Capture(packet, nil)

				rec.Logger.Printf("Event %d sent to sentry", eventID)
			}

			if rec.PrintStack {
				fmt.Fprintf(rw, f, err, stack)
			}
		}
	}()

	next(rw, r)
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf("Started %s %s", r.Method, r.URL.Path)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
