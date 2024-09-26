package GinTestContext

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"reflect"
)

var (
	ErrUnsupportedQueryType = errors.New("unsupported query type")
)

type queries struct {
	Queries interface{}
}

func (q *queries) SetQueries(queries interface{}) {
	q.Queries = queries
}

func (q *queries) writeQueriesToContext(c *gin.Context) error {
	if q.Queries == nil {
		return nil
	}

	switch reflect.ValueOf(q.Queries).Kind() {
	case reflect.Map:
		return q.writeQueriesWithMap(c)
	case reflect.Ptr, reflect.Struct:
		return q.writeQueriesWithObject(c)
	default:
		return ErrUnsupportedQueryType
	}
}

func (q *queries) writeQueriesWithMap(c *gin.Context) error {
	queries, ok := q.Queries.(map[string]string)
	if !ok {
		return ErrUnsupportedQueryType
	}

	values := url.Values{}
	for key, value := range queries {
		values.Set(key, value)
	}
	c.Request.URL = &url.URL{RawQuery: values.Encode()}

	return nil
}

func (q *queries) writeQueriesWithObject(c *gin.Context) error {
	var value reflect.Value
	switch reflect.ValueOf(q.Queries).Kind() {
	case reflect.Ptr:
		value = reflect.ValueOf(q.Queries).Elem()
		if value.Kind() != reflect.Struct {
			return ErrUnsupportedQueryType
		}
	case reflect.Struct:
		value = reflect.ValueOf(q.Queries)
	default:
		return ErrUnsupportedQueryType
	}

	values := url.Values{}
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}

		values.Set(tag, value.Field(i).String())
	}
	c.Request.URL = &url.URL{RawQuery: values.Encode()}

	return nil
}
