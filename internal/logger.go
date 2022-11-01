package internal

import (
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// RequestLogger logs every request
func RequestLogger(logger *zap.SugaredLogger, targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		//log request by who(IP address)
		requesterIP := r.RemoteAddr

		logger.Infow("Loaded page",
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
