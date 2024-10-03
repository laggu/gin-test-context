package GinTestContext

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Builder struct {
	headers
	uriParams
	queries
	body
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) GetContext() (*gin.Context, error) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{URL: &url.URL{}, Header: http.Header{}}

	if err := b.writeHeadersToContext(c); err != nil {
		return nil, err
	}

	if err := b.writeURIParamsToContext(c); err != nil {
		return nil, err
	}

	if err := b.writeQueriesToContext(c); err != nil {
		return nil, err
	}

	if err := b.writeBodyToContext(c); err != nil {
		return nil, err
	}

	return c, nil
}
