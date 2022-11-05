package handler

/*
import (
	"embed"
	"github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/sqlitestore"
	"github.com/koesie10/webauthn/webauthn"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Verify(t *testing.T) {
	mux := http.NewServeMux()
	type fields struct {
		config        config.Config
		logger        *zap.SugaredLogger
		statics       embed.FS
		sessionsStore *sqlitestore.SqliteStore
		webauth       *webauthn.WebAuthn
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   func() args
	}{
		{
			"happy path",
			fields{
				config:        config.Config{},
				logger:        zap.NewNop().Sugar(),
				statics:       embed.FS{},
				sessionsStore: nil,
				webauth:       nil,
			},
			func() args {
				req, _ := http.NewRequest("GET", "/verify", nil)
				return args{w: nil, r: req}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				config:        tt.fields.config,
				logger:        tt.fields.logger,
				statics:       tt.fields.statics,
				sessionsStore: tt.fields.sessionsStore,
				webauth:       tt.fields.webauth,
			}
			args := tt.args()
			h.Verify(args.w, args.r)
			executeRequest()
			//checkResponseCode(t, 200, )
		})
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
*/
