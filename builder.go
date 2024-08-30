package GinTestContext

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Builder struct {
	headers
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (m *Builder) GetContext() (*gin.Context, error) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{URL: &url.URL{}, Header: http.Header{}}

	if err := m.writeHeadersToContext(c); err != nil {
		return nil, err
	}

	return c, nil
}
