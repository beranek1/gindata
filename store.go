package gindata

import (
	"encoding/json"
	"strconv"

	"github.com/beranek1/godatainterface"
	"github.com/gin-gonic/gin"
)

type DataStoreErrorResponse struct {
	Error any
}

type DataStoreGetResponse struct {
	Data any
}

type DataStoreBackend struct {
	store godatainterface.DataStoreVersionedRangeFromInterval
}

func CreateDataStoreBackend(store godatainterface.DataStoreVersionedRangeFromInterval) *DataStoreBackend {
	return &DataStoreBackend{store: store}
}

func returnAsJsonResponse(c *gin.Context, data any, code int) {
	byt, err := json.Marshal(data)
	if err != nil {
		byt, err := json.Marshal(DataStoreGetResponse{err.Error()})
		if err != nil {
			c.String(500, "{\"Error\":\""+err.Error()+"\"}")
		} else {
			c.String(500, string(byt))
		}
	} else {
		c.String(code, string(byt))
	}
}

func (b *DataStoreBackend) Get(c *gin.Context) {
	key := c.Param("key")
	value, err := b.store.Get(key)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 404)
		return
	}
	returnAsJsonResponse(c, DataStoreGetResponse{value}, 200)
}

func (b *DataStoreBackend) GetAt(c *gin.Context) {
	key := c.Param("key")
	timestamp, err := strconv.ParseInt(c.Param("timestamp"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	value, err := b.store.GetAt(key, timestamp)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 404)
		return
	}
	returnAsJsonResponse(c, DataStoreGetResponse{value}, 200)
}

func (b *DataStoreBackend) Range(c *gin.Context) {
	key := c.Param("key")
	start, err := strconv.ParseInt(c.Param("start"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	end, err := strconv.ParseInt(c.Param("end"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	values, err := b.store.Range(key, start, end)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 404)
		return
	}
	returnAsJsonResponse(c, DataStoreGetResponse{values.Array()}, 200)
}

func (b *DataStoreBackend) From(c *gin.Context) {
	key := c.Param("key")
	start, err := strconv.ParseInt(c.Param("start"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	values, err := b.store.From(key, start)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 404)
		return
	}
	returnAsJsonResponse(c, DataStoreGetResponse{values.Array()}, 200)
}

func (b *DataStoreBackend) RangeInterval(c *gin.Context) {
	key := c.Param("key")
	start, err := strconv.ParseInt(c.Param("start"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	end, err := strconv.ParseInt(c.Param("end"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	interval, err := strconv.ParseInt(c.Param("interval"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	values, err := b.store.RangeInterval(key, start, end, interval)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 404)
		return
	}
	returnAsJsonResponse(c, DataStoreGetResponse{values.Array()}, 200)
}

func (b *DataStoreBackend) FromInterval(c *gin.Context) {
	key := c.Param("key")
	start, err := strconv.ParseInt(c.Param("start"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	interval, err := strconv.ParseInt(c.Param("interval"), 10, 64)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 400)
		return
	}
	values, err := b.store.FromInterval(key, start, interval)
	if err != nil {
		returnAsJsonResponse(c, DataStoreErrorResponse{err.Error()}, 404)
		return
	}
	returnAsJsonResponse(c, DataStoreGetResponse{values.Array()}, 200)
}

func (b *DataStoreBackend) AttachToRouter(path string, r *gin.Engine) {
	r.GET(path+"/:key", b.Get)
	r.GET(path+"/:key/at/:timestamp", b.GetAt)
	r.GET(path+"/:key/range/:start/:end/:interval", b.RangeInterval)
	r.GET(path+"/:key/range/:start/:end", b.Range)
	r.GET(path+"/:key/range/:start", b.From)
	r.GET(path+"/:key/from/:start", b.From)
	r.GET(path+"/:key/from/:start/:interval", b.FromInterval)
}

func (b *DataStoreBackend) SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	b.AttachToRouter("", r)
	return r
}
