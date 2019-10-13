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

// Logger struct to implement a http.HandlerFunc interface
type Logger struct {
	log       *log.Logger
	req       *http.Request
	status    int
	size      int
	timestamp time.Time
	res       http.ResponseWriter
}

// NewLogger returns a new Logger with provided writer prefix and flags
func NewLogger(out io.Writer, prefix string, flags int) *Logger {
	return &Logger{
		log: log.New(out, prefix, flags),
	}
}

// Header is used to get associated request Header details.
//Gets header from  ResponseWriter.Header
func (l *Logger) Header() http.Header {
	return l.res.Header()
}

// Write writes a header to response. Calls to ResponseWriter.Writer with provided byte slice
func (l *Logger) Write(b []byte) (int, error) {
	size, err := l.res.Write(b)
	l.size += size
	return size, err
}

// WriteHeader sends an HTTP response header with the provided status code to ResponseWriter.WriteHeader
func (l *Logger) WriteHeader(s int) {
	l.res.WriteHeader(s)
	l.status = s
}

// Status return the int status of Logger instance
func (l *Logger) Status() int {
	return l.status
}

// Size return the int size of Logger instance
func (l *Logger) Size() int {
	return l.size
}

// Flush calls the Flusher interface that is implemented by ResponseWriters
// that allow an HTTP handler to flush buffered data to the client.
func (l *Logger) Flush() {
	f, ok := l.res.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// Fatal logs fatal application errors and exits.
// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (l Logger) Fatal(v ...interface{}) {
	l.log.Fatal(v...)
}

// Fatalf logs fatal application errors using provided format and exits.
// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (l Logger) Fatalf(format string, v ...interface{}) {
	l.log.Fatalf(format, v...)
}

// Fatalln logs fatal application errors and exits.
// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func (l Logger) Fatalln(v ...interface{}) {
	l.log.Fatalln(v...)
}

// Panic prints error message passed and panics.
//Panic is equivalent to Print() followed by a call to panic().
func (l Logger) Panic(v ...interface{}) {
	l.log.Panic(v...)
}

// Panicf prints error message passed using provided format and panics.
// Panicf is equivalent to Printf() followed by a call to panic().
func (l Logger) Panicf(format string, v ...interface{}) {
	l.log.Panicf(format, v...)
}

// Panicln prints error message passed and panics.
// Panicln is equivalent to Println() followed by a call to panic().
func (l Logger) Panicln(v ...interface{}) {
	l.log.Panicln(v...)
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Println(v ...interface{}) {
	l.log.Println(v...)
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

// SetFlags sets the output flags for the Logger instance.
func (l Logger) SetFlags(flag int) {
	l.log.SetFlags(flag)
}

// SetOutput sets the output destination for the Logger instace.
func (l Logger) SetOutput(w io.Writer) {
	l.log.SetOutput(w)
}

// SetPrefix sets the output prefix for the Logger instace.
func (l Logger) SetPrefix(prefix string) {
	l.log.SetPrefix(prefix)
}

// Handler responds to an HTTP request.
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
