package GinTestContext

import (
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHeader struct {
	Foo    string `header:"foo"`
	Bar    string `header:"bar"`
	Ignore string
}

func Test_headers_SetHeaders(t *testing.T) {
	h := map[string]string{}
	err := faker.FakeData(&h)
	require.NoError(t, err)

	mock := headers{}
	mock.SetHeaders(h)
	require.Equal(t, h, mock.Headers)
}

func Test_headers_writeHeadersToContext(t *testing.T) {
	testCases := []struct {
		name           string
		header         interface{}
		expectedResult error
	}{
		{
			name:           "empty_header",
			header:         nil,
			expectedResult: nil,
		},
		{
			name:           "map",
			header:         map[string]string{},
			expectedResult: nil,
		},
		{
			name:           "object",
			header:         testHeader{},
			expectedResult: nil,
		},
		{
			name:           "object_pointer",
			header:         &testHeader{},
			expectedResult: nil,
		},
		{
			name:           "int",
			header:         1,
			expectedResult: ErrUnsupportedHeaderType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.header != nil {
				err := faker.FakeData(&tc.header)
				require.NoError(t, err)
			}

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = &http.Request{Header: http.Header{}}

			mock := headers{}
			mock.SetHeaders(tc.header)
			err := mock.writeHeadersToContext(c)
			require.Equal(t, tc.expectedResult, err)
		})
	}
}

func Test_headers_writeHeadersWithMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := map[string]string{}
		err := faker.FakeData(&h)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}

		mock := headers{}
		mock.SetHeaders(h)
		err = mock.writeHeadersWithMap(c)
		require.NoError(t, err)

		for key, value := range h {
			require.Equal(t, value, c.Request.Header.Get(key))
		}
	})
	t.Run("invalid_map_type", func(t *testing.T) {
		h := map[string]int{}
		err := faker.FakeData(&h)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}

		mock := headers{}
		mock.SetHeaders(h)
		err = mock.writeHeadersWithMap(c)
		require.Equal(t, ErrUnsupportedHeaderType, err)
	})
}

func Test_headers_writeHeadersWithObject(t *testing.T) {
	t.Run("success_with_pointer", func(t *testing.T) {
		header1 := &testHeader{}
		err := faker.FakeData(header1)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}

		mock := headers{}
		mock.SetHeaders(header1)
		err = mock.writeHeadersWithObject(c)
		require.NoError(t, err)

		header2 := &testHeader{}
		err = c.BindHeader(header2)
		require.NoError(t, err)
		require.Equal(t, header1.Foo, header2.Foo)
		require.Equal(t, header1.Bar, header2.Bar)
	})
	t.Run("success_with_object", func(t *testing.T) {
		header1 := testHeader{}
		err := faker.FakeData(&header1)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}

		mock := headers{}
		mock.SetHeaders(header1)
		err = mock.writeHeadersWithObject(c)
		require.NoError(t, err)

		header2 := &testHeader{}
		err = c.BindHeader(header2)
		require.NoError(t, err)
		require.Equal(t, header1.Foo, header2.Foo)
		require.Equal(t, header1.Bar, header2.Bar)
	})
	t.Run("fail_unsupported_pointer", func(t *testing.T) {
		var header int
		err := faker.FakeData(&header)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}

		mock := headers{}
		mock.SetHeaders(&header)
		err = mock.writeHeadersWithObject(c)
		require.Equal(t, ErrUnsupportedHeaderType, err)
	})
	t.Run("fail_unsupported_type", func(t *testing.T) {
		var header string
		err := faker.FakeData(&header)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}

		mock := headers{}
		mock.SetHeaders(header)
		err = mock.writeHeadersWithObject(c)
		require.Equal(t, ErrUnsupportedHeaderType, err)
	})
}
