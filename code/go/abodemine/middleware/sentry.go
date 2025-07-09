package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"
)

// SentryMiddlewareHandler wraps an httprouter.Handle
// to capture 5xx errors and send them to Sentry.
func SentryMiddlewareHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Create a Sentry hub for this request
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetRequest(r)
		ctx := sentry.SetHubOnContext(r.Context(), hub)

		// Start a transaction for this request
		transaction := sentry.StartTransaction(ctx, r.Method+" "+r.URL.Path)
		transaction.Status = sentry.HTTPtoSpanStatus(http.StatusOK)
		transaction.SetTag("http.method", r.Method)
		transaction.SetTag("http.url", r.URL.Path)
		transaction.SetTag("http.host", r.Host)

		// Add the transaction to the scope
		hub.Scope().SetSpan(transaction)

		// Start timing the request
		startTime := time.Now()

		rw := &ResposeWriter{ResponseWriter: w}

		// Call the next handler with our custom ResponseWriter.
		next(rw, r.WithContext(ctx), ps)

		// Calculate request duration.
		duration := time.Since(startTime)

		// Calculate the duration in milliseconds, using
		// float64 to avoid rounded results.
		durationMs := float64(duration.Microseconds()) / 1000.

		// Update transaction with response info
		transaction.Status = sentry.HTTPtoSpanStatus(rw.StatusCode)
		transaction.SetTag("http.status_code", strconv.Itoa(rw.StatusCode))
		transaction.SetData("http.duration_ms", durationMs)
		transaction.Finish()

		if rw.StatusCode < 500 {
			return
		}

		// Create a Sentry event for the error.

		event := sentry.NewEvent()
		event.Request = sentry.NewRequest(r)
		event.Transaction = r.URL.Path

		// Add status code to the event
		event.Extra["status_code"] = rw.StatusCode
		event.Extra["method"] = r.Method
		event.Extra["path"] = r.URL.Path
		event.Extra["duration_ms"] = durationMs

		// Create a custom error message.
		event.Message = http.StatusText(rw.StatusCode)

		// Capture the event
		hub.CaptureEvent(event)
	}
}
