package gm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"ops-server/global"
	"ops-server/model/system"
	"strings"
	"time"
)

type HttpClient struct {
	client  *http.Client
	headers map[string]string
	url     string
	param   map[string]string
}

type HttpResponse struct {
	Data interface{}
	Msg  string
	Code int
}

// NewHttpClient 创建一个新的 HTTP 客户端，允许自定义超时时间
func NewHttpClient(ctx context.Context, platform string) (client *HttpClient, err error) {
	var gamePlatform system.SysGamePlatform
	var gmUrl string

	if platform == "default" {
		gmUrl = global.OPS_CONFIG.Default.GmUrl
	} else {
		projectId := ctx.Value("projectId").(string)

		if projectId == "" {
			return
		}

		err = global.OPS_DB.First(&gamePlatform, "project_id = ? and platform_code = ?", projectId, platform).Error
		if err != nil {
			return
		}

		gmUrl = gamePlatform.GmUrl
	}

	if gmUrl == "" {
		return nil, errors.New("url is empty")
	} else if !strings.HasPrefix(gmUrl, "http") {
		return nil, errors.New("url is not http")
	}

	return &HttpClient{
		client: &http.Client{Timeout: time.Second * 30},
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		url: gmUrl,
	}, nil
}

// SetHeader 设置请求头
func (h *HttpClient) SetHeader(key, value string) {
	h.headers[key] = value
}

// request 通用请求方法
func (h *HttpClient) request(method, url string, body io.Reader) (response *HttpResponse, err error) {
	bashUrl := h.url + url
	req, err := http.NewRequest(method, bashUrl, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	result, err := h.client.Do(req)
	if err != nil {
		global.OPS_LOG.Error("request error", zap.Error(err))
		return nil, err
	} else if result.StatusCode != http.StatusOK {
		global.OPS_LOG.Error("request error", zap.Int("statusCode", result.StatusCode))
		return nil, errors.New("request error")
	}

	defer func() {
		err := result.Body.Close()
		if err != nil {
			global.OPS_LOG.Error("request error", zap.Error(err))
		}
	}()

	resultBody, err := io.ReadAll(result.Body)
	if err != nil {
		global.OPS_LOG.Error("request error", zap.Error(err))
		return
	}

	err = json.Unmarshal(resultBody, &response)
	if err != nil {
		global.OPS_LOG.Error("request error", zap.Error(err))
		return
	}
	global.OPS_LOG.Info("request result", zap.Any("response", response))

	if response.Code != 0 {
		global.OPS_LOG.Error("request error", zap.Any("response", response))
		return nil, errors.New(response.Msg)
	}
	return response, err
}
func (h *HttpClient) Get(uri string, params map[string]string) (*HttpResponse, error) {
	global.OPS_LOG.Info("Get", zap.String("baseURL", h.url+uri), zap.Any("params", params))
	// 构造查询参数
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	return h.request(http.MethodGet, u.String(), nil)
}

// Post 发送 POST 请求
func (h *HttpClient) Post(uri string, data []byte) (*HttpResponse, error) {
	global.OPS_LOG.Info("Post", zap.String("baseURL", h.url+uri), zap.ByteString("data", data))
	if h.param != nil {
		u, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}

		//query := u.Query()
		//for k, v := range h.param {
		//	query.Set(k, v)
		//}
		//u.RawQuery = query.Encode()
		var rawQuery []string
		for k, v := range h.param {
			rawQuery = append(rawQuery, k+"="+v)
		}
		u.RawQuery = strings.Join(rawQuery, "&")
		return h.request(http.MethodPost, u.String(), bytes.NewBuffer(data))
	}

	return h.request(http.MethodPost, uri, bytes.NewBuffer(data))
}
