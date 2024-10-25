package ginTestContext

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

type ContextInput struct {
	Headers   any
	URIParams any
	Queries   any
	Body      any
}

func NewTestContext(input ContextInput) (*gin.Context, error) {
	b := NewBuilder()
	b.SetHeaders(input.Headers)
	b.SetURIParams(input.URIParams)
	b.SetQueries(input.Queries)
	b.SetBody(input.Body)
	return b.GetContext()
}
