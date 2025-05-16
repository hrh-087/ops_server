package gm

import (
	"strings"
)

func (h HttpClient) DimensionPush(serverList []string) (response *HttpResponse, err error) {
	data := map[string]string{
		"serverIds": strings.Join(serverList, ","),
	}

	h.param = data

	response, err = h.Post("/dimension/push", []byte{})

	return
}
