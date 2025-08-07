package utils

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	validation "github.com/itgelo/ozzo-validation/v4"
	"net/http"
	"reflect"
	"time"
)

type Pagination struct {
	Records      int64 `json:"records"`
	TotalRecords int64 `json:"total_records"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	TotalPage    int   `json:"total_page"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ResponseWriter(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
	return
}

type data struct {
	Records interface{} `json:"records,omitempty"`
	Record  interface{} `json:"record,omitempty"`
}

type metaResponse struct {
	// CorrelationID is the response correlation_id
	//in: string
	CorrelationID string `json:"correlation_id"`
	// Code is the response code
	//in: int
	Code int `json:"code"`
	// Message is the response message
	//in: string
	Message string `json:"message"`
	//Time is the response message
	//in string
	Time string `json:"time"`
	// Pagination of the pagination response
	// in: PaginationResponse
	Pagination *Pagination `json:"pagination,omitempty"`
}

type responseHttp struct {
	// Meta is the API response information
	// in: MetaResponse
	Meta metaResponse `json:"meta"`
	// Data is our data
	// in: DataResponse
	Data data `json:"data"`
	// Errors is the response message
	//in: string
	Errors interface{} `json:"errors,omitempty"`
}

func setHttpResponse(code int, message string, result interface{}, paging *Pagination) interface{} {
	dt := data{}
	isSlice := reflect.ValueOf(result).Kind() == reflect.Slice
	if isSlice {
		dt.Records = result
		dt.Record = nil
	} else {
		dt.Records = nil
		dt.Record = result
	}

	correlationId := "req-" + uuid.NewString()
	current := time.Now().Format("2006-01-02 15:04:05")

	return responseHttp{
		Meta: metaResponse{
			CorrelationID: correlationId,
			Code:          code,
			Message:       message,
			Time:          current,
			Pagination:    paging,
		},
		Data: dt,
	}
}

func SetDefaultResponse(_ context.Context, msg Message) interface{} {
	return setHttpResponse(msg.Code, msg.Message, nil, nil)
}

func SetHttpResponse(_ context.Context, msg Message, result interface{}, paging *Pagination) interface{} {
	return setHttpResponse(msg.Code, msg.Message, result, paging)
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	msgResponse := Message{Code: http.StatusInternalServerError, Message: err.Error()}
	switch err.(type) {
	case validation.Errors:
		msgResponse = Message{Code: http.StatusBadRequest, Message: err.Error()}
	}
	ResponseWriter(w, http.StatusBadRequest, SetDefaultResponse(ctx, msgResponse))
}

type encodeError interface {
	error() error
}

func GetHttpResponse(resp interface{}) *responseHttp {
	if result, ok := resp.(responseHttp); ok {
		return &result
	}
	return nil
}

func EncodeResponseHTTP(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if err, ok := resp.(encodeError); ok && err.error() != nil {
		EncodeError(ctx, err.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	result := GetHttpResponse(resp)
	w.WriteHeader(result.Meta.Code)
	return json.NewEncoder(w).Encode(resp)
}
