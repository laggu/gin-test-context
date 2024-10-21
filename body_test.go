package ginTestContext

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func Test_body_SetBody(t *testing.T) {
	b := struct {
		A string
		B int
		C float64
	}{}
	err := faker.FakeData(&b)
	require.NoError(t, err)

	mock := body{}
	mock.SetBody(b)

	require.Equal(t, b, mock.Body)
}

func Test_body_writeBodyToContext(t *testing.T) {
	testCases := []struct {
		name           string
		body           interface{}
		expectedResult error
	}{
		{
			name: "success",
			body: struct {
				A string
				B int
				C float64
			}{
				A: "a",
				B: 1,
				C: 1.23,
			},
			expectedResult: nil,
		},
		{
			name:           "empty_body",
			body:           nil,
			expectedResult: nil,
		},
		{
			name:           "json_parsing_error",
			body:           make(chan int),
			expectedResult: ErrUnsupportedBodyType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := body{}
			mock.SetBody(tc.body)

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = &http.Request{}

			err := mock.writeBodyToContext(c)
			require.Equal(t, tc.expectedResult, err)
		})
	}
}
