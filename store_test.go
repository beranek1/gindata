package gindata

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beranek1/godata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type dataStoreMock struct {
	success      bool
	Last_call    string
	Return_value any
	Return_map   map[int64]any
}

func (m *dataStoreMock) Get(key string) (any, error) {
	m.Last_call = "Get " + key
	if m.success {
		return m.Return_value, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) GetAt(key string, timestamp int64) (any, error) {
	m.Last_call = "GetAt " + key + " " + fmt.Sprint(timestamp)
	if m.success {
		return m.Return_value, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) Put(key string, value any) error {
	m.Last_call = "Put " + key + " " + fmt.Sprint(value)
	if m.success {
		return nil
	}
	return errors.New("so sorry")
}

func (m *dataStoreMock) PutAt(key string, value any, timestamp int64) error {
	m.Last_call = "PutAt " + key + " " + fmt.Sprint(value) + " " + fmt.Sprint(timestamp)
	if m.success {
		return nil
	}
	return errors.New("so sorry")
}

func (m *dataStoreMock) Range(key string, start int64, end int64) (map[int64]any, error) {
	m.Last_call = "Range " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(end)
	if m.success {
		return m.Return_map, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) From(key string, start int64) (map[int64]any, error) {
	m.Last_call = "From " + key + " " + fmt.Sprint(start)
	if m.success {
		return m.Return_map, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) RangeInterval(key string, start int64, end int64, interval int64) (map[int64]any, error) {
	m.Last_call = "RangeInterval " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(end) + " " + fmt.Sprint(interval)
	if m.success {
		return m.Return_map, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) FromInterval(key string, start int64, interval int64) (map[int64]any, error) {
	m.Last_call = "FromInterval " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(interval)
	if m.success {
		return m.Return_map, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) RangeArray(key string, start int64, end int64) ([]godata.DataVersionArrayEntry, error) {
	m.Last_call = "Range " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(end)
	if m.success {
		return []godata.DataVersionArrayEntry{}, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) FromArray(key string, start int64) ([]godata.DataVersionArrayEntry, error) {
	m.Last_call = "From " + key + " " + fmt.Sprint(start)
	if m.success {
		return []godata.DataVersionArrayEntry{}, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) RangeIntervalArray(key string, start int64, end int64, interval int64) ([]godata.DataVersionArrayEntry, error) {
	m.Last_call = "RangeInterval " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(end) + " " + fmt.Sprint(interval)
	if m.success {
		return []godata.DataVersionArrayEntry{}, nil
	}
	return nil, errors.New("so sorry")
}

func (m *dataStoreMock) FromIntervalArray(key string, start int64, interval int64) ([]godata.DataVersionArrayEntry, error) {
	m.Last_call = "FromInterval " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(interval)
	if m.success {
		return []godata.DataVersionArrayEntry{}, nil
	}
	return nil, errors.New("so sorry")
}

func createTestBackendRouterSuccess() (*gin.Engine, *dataStoreMock) {
	rv := "value"
	rm := map[int64]any{}
	rm[123] = rv
	m := &dataStoreMock{true, "", rv, rm}
	b := CreateDataStoreBackend(m)
	return b.SetupRouter(), m
}

func createTestBackendRouterError() (*gin.Engine, *dataStoreMock) {
	m := &dataStoreMock{false, "", nil, nil}
	b := CreateDataStoreBackend(m)
	return b.SetupRouter(), m
}

func testPath(t *testing.T, method string, path string, code int, response string, call string) {
	router, m := createTestBackendRouterSuccess()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, code, w.Code)
	assert.Equal(t, response, w.Body.String())
	assert.Equal(t, m.Last_call, call)

	erouter, _ := createTestBackendRouterError()

	ew := httptest.NewRecorder()
	ereq, _ := http.NewRequest(method, path, nil)
	erouter.ServeHTTP(ew, ereq)

	assert.Equal(t, 404, ew.Code)
	assert.Equal(t, "{\"Error\":\"so sorry\"}", ew.Body.String())
}

func testIllegalPath(t *testing.T, method string, path string, code int, response string) {
	router, m := createTestBackendRouterSuccess()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, code, w.Code)
	assert.Equal(t, response, w.Body.String())
	assert.Equal(t, m.Last_call, "")
}

func TestGet(t *testing.T) {
	testPath(t, "GET", "/test", 200, "{\"Data\":\"value\"}", "Get test")
}

func TestGetAt(t *testing.T) {
	testPath(t, "GET", "/test/at/456", 200, "{\"Data\":\"value\"}", "GetAt test 456")
	testIllegalPath(t, "GET", "/test/at/abc", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}")
}

func TestRange(t *testing.T) {
	testPath(t, "GET", "/test/range/456", 200, "{\"Data\":{\"123\":\"value\"}}", "From test 456")
	testPath(t, "GET", "/test/range/456/789", 200, "{\"Data\":{\"123\":\"value\"}}", "Range test 456 789")
	testPath(t, "GET", "/test/range/456/789/10", 200, "{\"Data\":{\"123\":\"value\"}}", "RangeInterval test 456 789 10")
	testIllegalPath(t, "GET", "/test/range/abc", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/range/abc/789", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/range/456/def", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"def\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/range/abc/789/10", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/range/456/def/10", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"def\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/range/456/789/ghi", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"ghi\\\": invalid syntax\"}")
}

func TestFrom(t *testing.T) {
	testPath(t, "GET", "/test/from/456", 200, "{\"Data\":{\"123\":\"value\"}}", "From test 456")
	testPath(t, "GET", "/test/from/456/10", 200, "{\"Data\":{\"123\":\"value\"}}", "FromInterval test 456 10")
	testIllegalPath(t, "GET", "/test/from/abc", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/from/abc/789", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}")
	testIllegalPath(t, "GET", "/test/from/456/def", 400, "{\"Error\":\"strconv.ParseInt: parsing \\\"def\\\": invalid syntax\"}")
}
