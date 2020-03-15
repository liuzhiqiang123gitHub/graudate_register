package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"
)

// ReqData2Form try to parse request (from apiMiddleWare) body as json and inject user_id from header to body, if failed, deal with it as form.
// It should be called before your business logic.
func ReqData2Form() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("MiddleWare") == "ON" || c.Request.Header.Get("WM") == "ON" {
			reqData2Form(c)
		}
	}
}

// ReqData2Form try to parse all request body as json and inject user_id from header to body
func AllReqData2Form() gin.HandlerFunc {
	return reqData2Form
}

func reqData2Form(c *gin.Context) {
	userId := c.Request.Header.Get(CODOON_USER_ID)
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("read request body error:%v", err)
		return
	}
	// fmt.Printf("raw body:%s\n", data)
	var v map[string]interface{}
	if len(data) == 0 {
		v = make(map[string]interface{})
		err = nil
	} else {
		v, err = loadJson(bytes.NewReader(data))
	}
	if err != nil {
		// if request data is NOT json format, restore body
		// log.Printf("ReqData2Form parse as json failed. restore [%s] to body", string(data))
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
	} else {
		// if user_id in request is not empty, move it to req_user_id
		if uid, ok := v[CODOON_USER_ID]; ok {
			v["req_user_id"] = uid
		}
		// inject use_id into form
		v[CODOON_USER_ID] = userId
		form := map2Form(v)
		s := form.Encode()
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request.Body = ioutil.NopCloser(strings.NewReader(s))
		} else if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Request.Header.Del("Content-Type")
			// append url values
			urlValues := c.Request.URL.Query()
			for k, vv := range urlValues {
				if _, ok := form[k]; !ok {
					form[k] = vv
				}
			}
			c.Request.URL.RawQuery = form.Encode()
		} else {
			c.Request.Body = ioutil.NopCloser(strings.NewReader(s))
		}
	}
}

func loadJson(r io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var v map[string]interface{}
	err := decoder.Decode(&v)
	if err != nil {
		// log.Printf("loadJson decode error:%v", err)
		return nil, err
	}
	return v, nil
}

func map2Form(v map[string]interface{}) url.Values {
	form := url.Values{}
	var vStr string
	for key, value := range v {
		switch value.(type) {
		case string:
			vStr = value.(string)
		case float64, int, int64:
			vStr = fmt.Sprintf("%v", value)
		default:
			if b, err := json.Marshal(&value); err != nil {
				vStr = fmt.Sprintf("%v", value)
			} else {
				vStr = string(b)
			}
		}
		form.Set(key, vStr)
	}
	return form
}

//慢接口日志
type SlowLogger interface {
	Notice(format string, params ...interface{})
	Warning(format string, params ...interface{})
}

//慢接口日志
func GinSlowLogger(slog SlowLogger, threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()

		c.Next()

		endAt := time.Now()
		latency := endAt.Sub(startAt)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if latency > threshold {
			slog.Warning("[GIN Slowlog] %v | %3d | %12v | %s | %-7s %s %s\n%s",
				endAt.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				c.Request.URL.String(),
				c.Request.URL.Opaque,
				c.Errors.String())
		}
	}
}

const (
	CODOON_REQUEST_ID   = "codoon_request_id"
	CODOON_SERVICE_CODE = "codoon_service_code"
	CODOON_USER_ID      = "user_id"
	CODOON_DID          = "did"
	KAFKA_TOPIC         = "codoon-trace-log"
)
