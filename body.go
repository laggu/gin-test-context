package ginTestContext

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnsupportedBodyType = errors.New("unsupported body type")
)

type body struct {
	Body interface{}
}

func (b *body) SetBody(body interface{}) {
	b.Body = body
}

func (b *body) writeBodyToContext(c *gin.Context) error {
	if b.Body == nil {
		return nil
	}

	body, err := json.Marshal(b.Body)
	if err != nil {
		return ErrUnsupportedBodyType
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}
