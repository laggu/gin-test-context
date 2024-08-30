package GinTestContext

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnsupportedHeaderType = errors.New("unsupported header type")
)

type headers struct {
	Headers interface{}
}

func (h *headers) SetHeaders(headers interface{}) {
	h.Headers = headers
}

func (h *headers) writeHeadersToContext(c *gin.Context) error {
	if h.Headers == nil {
		return nil
	}

	switch reflect.ValueOf(h.Headers).Kind() {
	case reflect.Map:
		return h.writeHeadersWithMap(c)
	case reflect.Ptr, reflect.Struct:
		return h.writeHeadersWithObject(c)
	default:
		return ErrUnsupportedHeaderType
	}
}

func (h *headers) writeHeadersWithMap(c *gin.Context) error {
	headers, ok := h.Headers.(map[string]string)
	if !ok {
		return ErrUnsupportedHeaderType
	}
	for key, value := range headers {
		c.Request.Header.Add(key, value)
	}
	return nil
}

func (h *headers) writeHeadersWithObject(c *gin.Context) error {
	var value reflect.Value
	switch reflect.ValueOf(h.Headers).Kind() {
	case reflect.Ptr:
		value = reflect.ValueOf(h.Headers).Elem()
		if value.Kind() != reflect.Struct {
			return ErrUnsupportedHeaderType
		}
	case reflect.Struct:
		value = reflect.ValueOf(h.Headers)
	default:
		return ErrUnsupportedHeaderType
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("header")
		if tag == "" {
			continue
		}
		c.Request.Header.Add(tag, value.Field(i).String())
	}

	return nil
}
