package notice

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type DingTalk struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
}

// SendDingTalkMessage 发送通用钉钉消息
func SendDingTalkMessage(webhook, secret string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	// 处理加签
	if secret != "" {
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		sign, err := signDingTalk(secret, timestamp)
		if err != nil {
			return err
		}
		webhook += fmt.Sprintf("&timestamp=%s&sign=%s", timestamp, sign)
	}

	// 发送请求
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("post error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dingtalk webhook failed, status: %d", resp.StatusCode)
	}

	return nil
}

// signDingTalk 加签算法
func signDingTalk(secret, timestamp string) (string, error) {
	stringToSign := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(sign), nil
}
