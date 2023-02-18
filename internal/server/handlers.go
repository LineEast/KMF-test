package server

import (
	"strings"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

type (
	Request struct {
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers"`
		Url     string            `json:"url"`
		Body    []byte            `json:"body"`
	}

	Response struct {
		ID      int               `json:"id"`
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Length  int               `json:"length"`
	}
)

func (s *Server) Main() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Получаем запрос ===
		request := Request{}

		err := json.Unmarshal(ctx.Request.Body(), &request)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		// Делаем запрос
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod(request.Method)
		req.SetRequestURI(request.Url)

		for i := range request.Headers {
			req.Header.Set(i, request.Headers[i])
		}

		if request.Method == fasthttp.MethodPost && len(request.Body) != 0 {
			req.SetBody(request.Body)
		}

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		client := &fasthttp.Client{}
		err = client.Do(req, resp)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		// Записываем резултаты в базу и получаем ID
		a := Request{
			Method:  string(ctx.Request.Header.Method()),
			Headers: makeHeadersMap(string(ctx.Request.Header.Header())),
			Url:     ctx.Request.URI().String(),
			Body:    ctx.Request.Body(),
		}

		r, err := json.Marshal(&a)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		c := Request{
			Method:  request.Method,
			Headers: makeHeadersMap(resp.Header.String()),
			Body:    resp.Body(),
		}

		rr, err := json.Marshal(&c)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		id, err := s.database.Create(ctx, r, rr)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		// Подготавливаем и отдаем ответ клиенту
		headers := map[string]string{}
		resp.Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})

		response := Response{
			ID:      id,
			Status:  resp.StatusCode(),
			Length:  resp.Header.ContentLength(),
			Headers: headers,
		}

		b, err := json.Marshal(response)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBody(b)
	}
}

func makeHeadersMap(headers string) map[string]string {
	first := strings.Split(headers, "\n")
	m := make(map[string]string)
	for i := range first {
		second := strings.Split(first[i], ":")
		if len(second) == 2 {
			m[second[0]] = second[1]
		}
	}

	return m
}
