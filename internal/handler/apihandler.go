package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/repository"
)

func AddShorten(resp http.ResponseWriter, req *http.Request) {

	var body []byte
	var err error
	var request APIRequest
	var response APIResponse

	defer req.Body.Close()
	if body, err = io.ReadAll(req.Body); err != nil || len(body) == 0 {
		response = APIResponse{Error: `url required`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		response = APIResponse{Error: `url required`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	if request.URL == `` {
		response = APIResponse{Error: `url required`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	if _, err = url.ParseRequestURI(request.URL); err != nil {
		response = APIResponse{Error: `Url invalid`, StatusCode: http.StatusBadRequest}
		sendAPIResponse(resp, &response)
		return
	}

	response = APIResponse{
		Result:     config.GetOptions().BaseHost + `/` + repository.GetRepository().Add(request.URL),
		StatusCode: http.StatusOK,
	}

	sendAPIResponse(resp, &response)
}

func sendAPIResponse(respWr http.ResponseWriter, resp *APIResponse) {
	rawByte, err := json.Marshal(resp)
	if err != nil {
		logger.Error(`response marshal error`, err)
		respWr.WriteHeader(http.StatusInternalServerError)
		return
	}

	respWr.Header().Add(HeaderContentTypeName, HeaderContentTypeJSONValue)
	respWr.WriteHeader(resp.StatusCode)

	_, err = respWr.Write(rawByte)
	if err != nil {
		logger.Error(`write response error`, err)
		respWr.WriteHeader(http.StatusInternalServerError)
		return
	}
}
