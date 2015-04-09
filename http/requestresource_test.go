package http_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/patdowney/downloaderd-request/api"
	dh "github.com/patdowney/downloaderd-request/http"
)

func TestURLResolving(t *testing.T) {
	requestID := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

	expectedURL, _ := url.Parse("/request/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	res := dh.RequestResource{}

	router := mux.NewRouter()

	res.RegisterRoutes(router.PathPrefix("/request").Subrouter())

	u, err := res.GetRequestURL(requestID)

	if err != nil {
		t.Error(err)
	}

	if *u != *expectedURL {
		t.Errorf(`GetRequestURL('%s') = %q want %q`, requestID, u, expectedURL)
	}
}

func TestGetRequestURL(t *testing.T) {
	requestID := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

	expectedURL, _ := url.Parse("/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	res := dh.RequestResource{}
	router := mux.NewRouter()

	res.RegisterRoutes(router)

	u, err := res.GetRequestURL(requestID)

	if err != nil {
		t.Error(err)
	}

	if *u != *expectedURL {
		t.Errorf(`GetRequestURL('%s') = %q want %q`, requestID, u, expectedURL)
	}
}

func TestParseCompleteRequest(t *testing.T) {
	res := dh.RequestResource{}

	jsonString := `{ "url": "http://example.com/some/resource", "checksum": "abcde", "checksum_type": "abc", "callback": "http://example.com/callback" }`
	incomingJSON := strings.NewReader(jsonString)

	expectedIncoming := api.IncomingRequest{
		URL:          "http://example.com/some/resource",
		Checksum:     "abcde",
		ChecksumType: "abc",
		Callback:     "http://example.com/callback"}

	r, _ := res.DecodeInputRequest(incomingJSON)
	if *r != expectedIncoming {
		t.Errorf(`DecodeInputRequest('%s') = %q want %q`, jsonString, r, expectedIncoming)
	}
}

func TestParseWithPartialRequest(t *testing.T) {
	res := dh.RequestResource{}

	jsonString := `{ "url": "http://example.com/some/resource" }`
	incomingJSON := strings.NewReader(jsonString)

	expectedIncoming := api.IncomingRequest{URL: "http://example.com/some/resource"}

	r, _ := res.DecodeInputRequest(incomingJSON)
	if *r != expectedIncoming {
		t.Errorf(`DecodeInputRequest('%s') = %q want %q`, jsonString, r, expectedIncoming)
	}
}

func TestRequestResourceGetIndex(t *testing.T)            {}
func TestRequestResourceGetRequest(t *testing.T)          {}
func TestRequestResourcePostIncomingRequest(t *testing.T) {}

//
//func (r *RequestResource) Index() HandlerFunc {
//	return func(rw http.ResponseWriter, req *http.Request) {
//		rw.Header().Set("Content-Type", "text/json")
//		rw.WriteHeader(http.StatusOK)
//		encoder := json.NewEncoder(rw)
//		encoder.Encode(r.RequestStore.ListAll())
//	}
//}
//
//func (r *RequestResource) Get() HandlerFunc {
//	return func(rw http.ResponseWriter, req *http.Request) {
//		vars := mux.Vars(req)
//		requestID := vars["requestID"]
//
//		downloadRequest := r.RequestStore.FindByID(requestID)
//
//		if downloadRequest != nil {
//			rw.Header().Set("Content-Type", "text/json")
//			rw.WriteHeader(http.StatusOK)
//			encoder := json.NewEncoder(rw)
//			encoder.Encode(downloadRequest)
//		} else {
//			rw.WriteHeader(http.StatusNotFound)
//			log.Printf("Couldn't find request with id:%s", downloadRequestID)
//		}
//	}
//}
//
//func (r *RequestResource) decodeIncomingRequest(body io.Reader) (*ApiIncomingRequest, error) {
//	decoder := json.NewDecoder(body)
//	var incomingRequest ApiIncomingRequest
//	err := decoder.Decode(&incomingRequest)
//	if err != nil {
//		return nil, err
//	}
//
//	_, err = incomingRequest.IsValid()
//	if err != nil {
//		return nil, err
//	}
//
//	return &incomingRequest
//}
//
//func (r *RequestResource) Post() HandlerFunc {
//	return func(http.ResponseWriter, *http.Request) {
//		incomingRequest, err := r.decodeInputRequest(req.Body)
//		if err != nil {
//			http.Error(rw, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		d, err := downloader.ProcessIncomingRequest(&incomingRequest)
//
//		if err != nil {
//			rw.WriteHeader(http.StatusInternalServerError)
//			log.Printf("server-error: %v", err)
//		} else {
//			newURL := GetURLForDownloadRequestID(req, d.ID)
//
//			rw.Header().Set("Content-Type", "text/json")
//			rw.Header().Set("Location", newURL)
//			rw.WriteHeader(http.StatusAccepted)
//			encoder := json.NewEncoder(rw)
//			encoder.Encode(d)
//
//		}
//	}
//}
