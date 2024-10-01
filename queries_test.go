package GinTestContext

import (
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testQueries struct {
	Foo    string `form:"foo"`
	Bar    string `form:"bar"`
	Ignore string
}

func Test_queries_SetQueries(t *testing.T) {
	q := testQueries{}
	err := faker.FakeData(&q)
	require.NoError(t, err)

	mock := queries{}
	mock.SetQueries(q)
	require.Equal(t, q, mock.Queries)
}

func Test_queries_writeQueriesToContext(t *testing.T) {
	testCases := []struct {
		name           string
		queries        interface{}
		expectedResult error
	}{
		{
			name:           "empty",
			queries:        nil,
			expectedResult: nil,
		},
		{
			name:           "map",
			queries:        map[string]string{},
			expectedResult: nil,
		},
		{
			name:           "object",
			queries:        testHeader{},
			expectedResult: nil,
		},
		{
			name:           "object_pointer",
			queries:        &testHeader{},
			expectedResult: nil,
		},
		{
			name:           "int",
			queries:        1,
			expectedResult: ErrUnsupportedQueryType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.queries != nil {
				err := faker.FakeData(&tc.queries)
				require.NoError(t, err)
			}

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = &http.Request{}

			mock := queries{}
			mock.SetQueries(tc.queries)
			err := mock.writeQueriesToContext(c)
			require.Equal(t, tc.expectedResult, err)
		})
	}
}

func Test_queries_writeQueriesWithMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := map[string]string{}
		err := faker.FakeData(&q)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := queries{}
		mock.SetQueries(q)
		err = mock.writeQueriesWithMap(c)
		require.NoError(t, err)

		for key, value := range q {
			require.Equal(t, value, c.Request.URL.Query().Get(key))
		}
	})
	t.Run("invalid_map_type", func(t *testing.T) {
		q := map[string]int{}
		err := faker.FakeData(&q)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := queries{}
		mock.SetQueries(q)
		err = mock.writeQueriesWithMap(c)
		require.Equal(t, ErrUnsupportedQueryType, err)
	})
}

func Test_queries_writeQueriesWithObject(t *testing.T) {
	t.Run("success_with_pointer", func(t *testing.T) {
		queries1 := &testQueries{}
		err := faker.FakeData(queries1)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := queries{}
		mock.SetQueries(queries1)
		err = mock.writeQueriesWithObject(c)
		require.NoError(t, err)

		queries2 := &testQueries{}
		err = c.BindQuery(queries2)
		require.NoError(t, err)
		require.Equal(t, queries1.Foo, queries2.Foo)
		require.Equal(t, queries1.Bar, queries2.Bar)
	})
	t.Run("success_with_object", func(t *testing.T) {
		queries1 := testQueries{}
		err := faker.FakeData(&queries1)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := queries{}
		mock.SetQueries(queries1)
		err = mock.writeQueriesWithObject(c)
		require.NoError(t, err)

		queries2 := &testQueries{}
		err = c.BindQuery(queries2)
		require.NoError(t, err)
		require.Equal(t, queries1.Foo, queries2.Foo)
		require.Equal(t, queries1.Bar, queries2.Bar)
	})
	t.Run("fail_unsupported_pointer", func(t *testing.T) {
		var q int
		err := faker.FakeData(&q)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := queries{}
		mock.SetQueries(&q)
		err = mock.writeQueriesWithObject(c)
		require.Equal(t, ErrUnsupportedQueryType, err)
	})
	t.Run("fail_unsupported_type", func(t *testing.T) {
		var q string
		err := faker.FakeData(&q)
		require.NoError(t, err)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{}

		mock := queries{}
		mock.SetQueries(q)
		err = mock.writeQueriesWithObject(c)
		require.Equal(t, ErrUnsupportedQueryType, err)
	})
}
