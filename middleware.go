package goezyrouting

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

type middleware func(next http.Handler) http.Handler

var prefix string
var reqid uint64

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			myid := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%s-%06d", prefix, myid)
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestURLPath, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func Recoverer(l Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if l != nil {
						l.Panic("app failed", rvr, debug.Stack())
					} else {
						fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
						debug.PrintStack()
					}
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func logger(l Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			t1 := time.Now()
			defer func() {
				l.Log("Served:",
					"proto", r.Proto,
					"method", r.Method,
					"path", path,
					"took", time.Since(t1),
					"reqId", getReqID(r.Context()),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headers := w.Header()
		origin := r.Header.Get("Origin")

		// Always set Vary headers
		// see https://github.com/rs/cors/issues/10,
		//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		headers.Set("Access-Control-Allow-Origin", origin)
		// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
		// by Access-Control-Request-Method (if supported) can be enough
		headers.Set("Access-Control-Allow-Methods", strings.ToUpper(r.Header.Get("Access-Control-Request-Method")))

		headers.Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Protects from MimeType Sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Prevents browser from prefetching DNS
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		// Denies website content to be served in an iframe
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
		// Prevents Internet Explorer from executing downloads in site's context
		w.Header().Set("X-Download-Options", "noopen")
		// Minimal XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I'm the first API middleware... Maybe do some cors stuff")
		next.ServeHTTP(w, r)
	})
}
func apiMiddleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I'm the second API middleware... Maybe do some authentication")
		next.ServeHTTP(w, r)
	})
}

func requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("This is a required path. If a user isn't in context maybe i'll logout")
		next.ServeHTTP(w, r)
	})
}
func minAccessLevel(minAccessLvl int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("User needs min access level:", minAccessLvl)
			next.ServeHTTP(w, r)
		})
	}
}
func levelTwoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("level two middleware")
		next.ServeHTTP(w, r)
	})
}
