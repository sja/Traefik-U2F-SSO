package internal

import (
	"errors"
	"fmt"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"syscall"
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

func InitLogger() (*zap.SugaredLogger, error) {
	zaplog, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	if viper.GetBool(ConfDebug) {
		zaplog, _ = zap.NewDevelopment()
	}
	logger := zaplog.Sugar()

	if err := logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		// see https://github.com/uber-go/zap/issues/880
		panic(fmt.Errorf("fatal error initializing logger: %w", err))
	}
	return logger, nil
}
