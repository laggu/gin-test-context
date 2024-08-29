package GinTestContext

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type Builder struct {
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (m *Builder) GetContext() (*gin.Context, error) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{URL: &url.URL{}}

	return c, nil
}
