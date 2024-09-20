package GinTestContext

import (
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testUriParam struct {
	Foo    string `uri:"foo"`
	Bar    string `uri:"bar"`
	Ignore string
}

func Test_uriParams_SetURIParams(t *testing.T) {
	p := testUriParam{}
	err := faker.FakeData(&p)
	require.NoError(t, err)

	mock := uriParams{}
	mock.SetURIParams(p)
	require.Equal(t, p, mock.URIParams)
}

func Test_uriParams_writeURIParamsToContext(t *testing.T) {
	testCases := []struct {
		name           string
		uriParam       interface{}
		expectedResult error
	}{
		{
			name:           "empty_uri",
			uriParam:       nil,
			expectedResult: nil,
		},
		{
			name:           "map",
			uriParam:       map[string]string{},
			expectedResult: nil,
		},
		{
			name:           "object",
			uriParam:       testHeader{},
			expectedResult: nil,
		},
		{
			name:           "object_pointer",
			uriParam:       &testHeader{},
			expectedResult: nil,
		},
		{
			name:           "int",
			uriParam:       1,
			expectedResult: ErrUnsupportedURIParamType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.uriParam != nil {
				err := faker.FakeData(&tc.uriParam)
				require.NoError(t, err)
			}

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = &http.Request{}

			mock := uriParams{}
			mock.SetURIParams(tc.uriParam)
			err := mock.writeURIParamsToContext(c)
			require.Equal(t, tc.expectedResult, err)
		})
	}
}

func Test_uriParams_writeWithMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := map[string]string{}
		err := faker.FakeData(&p)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := uriParams{}
		mock.SetURIParams(p)
		err = mock.writeURIParamsWithMap(c)
		require.NoError(t, err)

		require.Equal(t, len(c.Params), len(p))
		for _, param := range c.Params {
			require.Equal(t, p[param.Key], param.Value)
		}
	})
	t.Run("invalid_map_type", func(t *testing.T) {
		p := map[string]int{}
		err := faker.FakeData(&p)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := uriParams{}
		mock.SetURIParams(p)
		err = mock.writeURIParamsWithMap(c)
		require.Equal(t, ErrUnsupportedURIParamType, err)
	})
}

func Test_uriParams_writeWithObject(t *testing.T) {
	t.Run("success_with_pointer", func(t *testing.T) {
		uriParam1 := &testUriParam{}
		err := faker.FakeData(uriParam1)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := uriParams{}
		mock.SetURIParams(uriParam1)
		err = mock.writeURIParamsWithObject(c)
		require.NoError(t, err)

		uriParam2 := &testUriParam{}
		err = c.BindUri(uriParam2)
		require.NoError(t, err)
		require.Equal(t, uriParam1.Foo, uriParam2.Foo)
		require.Equal(t, uriParam1.Bar, uriParam2.Bar)
	})
	t.Run("success_with_object", func(t *testing.T) {
		uriParam1 := testUriParam{}
		err := faker.FakeData(&uriParam1)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := uriParams{}
		mock.SetURIParams(uriParam1)
		err = mock.writeURIParamsWithObject(c)
		require.NoError(t, err)

		uriParam2 := &testUriParam{}
		err = c.BindUri(uriParam2)
		require.NoError(t, err)
		require.Equal(t, uriParam1.Foo, uriParam2.Foo)
		require.Equal(t, uriParam1.Bar, uriParam2.Bar)
	})
	t.Run("fail_unsupported_pointer", func(t *testing.T) {
		var p int
		err := faker.FakeData(&p)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := uriParams{}
		mock.SetURIParams(&p)
		err = mock.writeURIParamsWithObject(c)
		require.Equal(t, ErrUnsupportedURIParamType, err)
	})
	t.Run("fail_unsupported_type", func(t *testing.T) {
		var p string
		err := faker.FakeData(&p)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := uriParams{}
		mock.SetURIParams(p)
		err = mock.writeURIParamsWithObject(c)
		require.Equal(t, ErrUnsupportedURIParamType, err)
	})
}
