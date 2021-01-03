package router

import (
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"project-name/config"
	"project-name/pkg/xlogger"
	"testing"
)

func TestInitRoute(t *testing.T) {
	config.InitArgs()
	err := config.InitConfig(config.Arg.Env, "../conf")
	if err != nil {
		t.Error(err)
	}
	xlogger.InitLogger()

	g := InitRoute()
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/sd/health", nil)
	g.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}
