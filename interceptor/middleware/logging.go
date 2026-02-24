package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"

	"github.com/kedacore/http-add-on/pkg/util"
)

const (
	CombinedLogFormat     = `%s %s %s [%s] "%s %s %s" %d %d "%s" "%s"`
	CombinedLogTimeFormat = "02/Jan/2006:15:04:05 -0700"
	CombinedLogBlankValue = "-"
)

type Logging struct {
	logger          logr.Logger
	upstreamHandler http.Handler
}

func NewLogging(logger logr.Logger, upstreamHandler http.Handler) *Logging {
	return &Logging{
		logger:          logger,
		upstreamHandler: upstreamHandler,
	}
}

var _ http.Handler = (*Logging)(nil)

func (lm *Logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = util.RequestWithLogger(r, lm.logger.WithName("LoggingMiddleware"))
	w = newResponseWriter(w)

	startTime := time.Now()
	defer lm.logAsync(w, r, startTime)

	lm.upstreamHandler.ServeHTTP(w, r)
}

func (lm *Logging) logAsync(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	go lm.log(w, r, startTime)
}

func (lm *Logging) log(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	ctx := r.Context()
	logger := util.LoggerFromContext(ctx)

	lrw := w.(*responseWriter)
	if lrw == nil {
		lrw = newResponseWriter(w)
	}

	timestamp := startTime.Format(CombinedLogTimeFormat)
	log := fmt.Sprintf(
		CombinedLogFormat,
		r.RemoteAddr,
		CombinedLogBlankValue,
		CombinedLogBlankValue,
		timestamp,
		r.Method,
		r.URL.Path,
		r.Proto,
		lrw.StatusCode(),
		lrw.BytesWritten(),
		r.Referer(),
		r.UserAgent(),
	)
	logger.Info(log)
}
