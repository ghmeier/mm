package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ghmeier/go-service"
)

type Handler interface {
	Init(string)
	Go([]string)
	Cmd() string
	Help()
}

type base struct {
	s   service.Service
	cmd string
	url string
	key string
}

func new(key, cmd, path string) *base {
	return &base{
		cmd: cmd,
		key: key,
		url: fmt.Sprintf("https://api.mixmax.com/v1/%s", path),
		s:   service.NewCustom(&responder{}),
	}
}

func (b *base) Cmd() string {
	return b.cmd
}

func (b *base) Key() string {
	return b.key
}

func (b *base) Url() string {
	return b.url
}

type responder struct {
}

type Response struct {
	Message string `json: "message,omitempty"`
	body    []byte
}

func (r *responder) Marshal(res *http.Response) (service.Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var raw Response
	if body == nil || len(body) == 0 {
		raw.body = make([]byte, 0)
		return &raw, nil
	}

	if strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		raw.Message = string(body)
		return &raw, nil

	}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	raw.body = body
	return &raw, nil
}

func (r *Response) Error() error {
	if r.Message == "" {
		return nil
	}

	return fmt.Errorf(r.Message)
}

func (r *Response) Body() ([]byte, error) {
	return r.body, nil
}
