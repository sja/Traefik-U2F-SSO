package internal

import (
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// RequestLogger logs every request
func RequestLogger(logger *zap.SugaredLogger, targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		//log requester ip
		requesterIP := r.RemoteAddr
		fwdAddress := r.Header.Get("x-forwarded-for")

		if fwdAddress != "" {
			// Got X-Forwarded-For
			requesterIP = fwdAddress // If it's a single IP, then awesome!

			// If we got an array... grab the first IP
			ips := strings.Split(fwdAddress, ", ")
			if len(ips) > 1 {
				requesterIP = ips[0]
			}
		}

		logger.Infow("handled request",
			"Method", r.Method,
			"RequestURI", r.RequestURI,
			"RequesterIP", requesterIP,
			"Time", time.Since(start),
		)
	})
}

func InitLogger(config Config) (*zap.SugaredLogger, error) {
	var zapcfg zap.Config
	if config.Debug {
		zapcfg = zap.NewDevelopmentConfig()
	} else {
		zapcfg = zap.NewProductionConfig()
	}
	zaplog, err := zapcfg.Build()
	if err != nil {
		return nil, err
	}

	logger := zaplog.Sugar()
	return logger, nil
}
