package GinTestContext

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
)

var (
	ErrUnsupportedURIParamType = errors.New("unsupported uri param type")
)

type uriParams struct {
	URIParams interface{}
}

func (u *uriParams) SetURIParams(uriParams interface{}) {
	u.URIParams = uriParams
}

func (u *uriParams) writeURIParamsToContext(c *gin.Context) error {
	if u.URIParams == nil {
		return nil
	}

	switch reflect.ValueOf(u.URIParams).Kind() {
	case reflect.Map:
		return u.writeURIParamsWithMap(c)
	case reflect.Ptr, reflect.Struct:
		return u.writeURIParamsWithObject(c)
	default:
		return ErrUnsupportedURIParamType
	}
}

func (u *uriParams) writeURIParamsWithMap(c *gin.Context) error {
	uriParams, ok := u.URIParams.(map[string]string)
	if !ok {
		return ErrUnsupportedURIParamType
	}
	for key, value := range uriParams {
		c.Params = append(c.Params, gin.Param{Key: key, Value: value})
	}
	return nil
}

func (u *uriParams) writeURIParamsWithObject(c *gin.Context) error {
	var value reflect.Value
	switch reflect.ValueOf(u.URIParams).Kind() {
	case reflect.Ptr:
		value = reflect.ValueOf(u.URIParams).Elem()
		if value.Kind() != reflect.Struct {
			return ErrUnsupportedURIParamType
		}
	case reflect.Struct:
		value = reflect.ValueOf(u.URIParams)
	default:
		return ErrUnsupportedURIParamType
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("uri")
		if tag == "" {
			continue
		}
		c.Params = append(c.Params, gin.Param{Key: tag, Value: value.Field(i).String()})
	}

	return nil
}
