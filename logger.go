package logger

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var now = func() time.Time {
	return time.Now().UTC()
}

type Logger struct {
	log       *log.Logger
	req       *http.Request
	status    int
	size      int
	timestamp time.Time
	res       http.ResponseWriter
}

func NewLogger(out io.Writer, prefix string, flags int) *Logger {
	return &Logger{
		log: log.New(out, prefix, flags),
	}
}

func (l *Logger) Header() http.Header {
	return l.res.Header()
}

func (l *Logger) Write(b []byte) (int, error) {
	size, err := l.res.Write(b)
	l.size += size
	return size, err
}

func (l *Logger) WriteHeader(s int) {
	l.res.WriteHeader(s)
	l.status = s
}

func (l *Logger) Status() int {
	return l.status
}

func (l *Logger) Size() int {
	return l.size
}

func (l *Logger) Flush() {
	f, ok := l.res.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func (l Logger) Fatal(v ...interface{}) {
	l.log.Fatal(v...)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.log.Fatalf(format, v...)
}

func (l Logger) Fatalln(v ...interface{}) {
	l.log.Fatalln(v...)
}

func (l Logger) Panic(v ...interface{}) {
	l.log.Panic(v...)
}

func (l Logger) Panicf(format string, v ...interface{}) {
	l.log.Panicf(format, v...)
}

func (l Logger) Panicln(v ...interface{}) {
	l.log.Panicln(v...)
}

func (l Logger) Println(v ...interface{}) {
	l.log.Println(v...)
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

func (l *Logger) Handler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.res = w
		l.req = r
		l.status = http.StatusOK
		l.timestamp = now()

		handler.ServeHTTP(l, r)
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}

		f := now()
		l.Printf(
			"%s - - [%s] \"%s %s %s\" %d %d %q %q %.4f\n",
			l.host(),
			f.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			l.uri(),
			r.Proto,
			l.Status(),
			l.Size(),
			l.referer(),
			l.userAgent(),
			f.Sub(l.timestamp).Seconds(),
		)
	}
}

func (l Logger) host() (host string) {
	var err error
	if host, _, err = net.SplitHostPort(l.req.RemoteAddr); err != nil {
		host = l.req.RemoteAddr
	}
	return
}

func (l Logger) uri() (uri string) {
	uri = l.req.RequestURI
	if l.req.ProtoMajor == 2 && l.req.Method == "CONNECT" {
		uri = l.req.Host
	}
	if uri == "" {
		uri = l.req.URL.RequestURI()
	}
	return
}

func (l Logger) referer() (referer string) {
	referer = l.req.Referer()
	if referer == "" {
		referer = "-"
	}
	return
}

func (l Logger) userAgent() (userAgent string) {
	userAgent = l.req.UserAgent()
	if userAgent == "" {
		userAgent = "-"
	}
	return
}
