package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/patdowney/downloaderd-request/api"
	"github.com/patdowney/downloaderd-request/common"
	"github.com/patdowney/downloaderd-request/download"
)

// RequestResource ...
type RequestResource struct {
	Clock          common.Clock
	RequestService *download.RequestService
	router         *mux.Router
	linkResolver   *api.LinkResolver
}

// NewRequestResource ...
func NewRequestResource(requestService *download.RequestService, linkResolver *api.LinkResolver) *RequestResource {
	return &RequestResource{
		Clock:          &common.RealClock{},
		RequestService: requestService,
		linkResolver:   linkResolver}
}

// RegisterRoutes ...
func (r *RequestResource) RegisterRoutes(parentRouter *mux.Router) {
	parentRouter.HandleFunc("/", r.Index()).Methods("GET", "HEAD")
	parentRouter.HandleFunc("/", r.Post()).Methods("POST")
	// regexp matches ids that look like '8671301b-49fa-416c-4bc0-2869963779e5'
	parentRouter.HandleFunc("/{id:[a-f0-9-]{36}}", r.Get()).Methods("GET", "HEAD").Name("request")

	r.router = parentRouter
}

func (r *RequestResource) populateListLinks(req *http.Request, requestList *[]*api.Request) {
	for _, l := range *requestList {
		r.populateLinks(req, l)
	}
}

func (r *RequestResource) populateLinks(req *http.Request, request *api.Request) {
	request.ResolveLinks(r.linkResolver, req)
}

// WrapError ...
func (r *RequestResource) WrapError(err error) *api.Error {
	return download.ToAPIError(common.NewErrorWrapper(err, r.Clock.Now()))
}

// Index ...
func (r *RequestResource) Index() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		requestList, err := r.RequestService.ListAll()

		encoder := json.NewEncoder(rw)
		rw.Header().Set("Content-Type", "application/json")

		if err != nil {
			log.Printf("server-error: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			encErr := encoder.Encode(r.WrapError(err))
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}
		} else {
			rw.WriteHeader(http.StatusOK)
			rl := download.ToAPIRequestList(&requestList)
			r.populateListLinks(req, rl)
			encErr := encoder.Encode(rl)
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}
		}
	}
}

// Get ...
func (r *RequestResource) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		requestID := vars["id"]

		downloadRequest, err := r.RequestService.FindByID(requestID)

		encoder := json.NewEncoder(rw)
		rw.Header().Set("Content-Type", "application/json")

		if err != nil {
			log.Printf("server-error: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			encErr := encoder.Encode(r.WrapError(err))
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}
		} else if downloadRequest != nil {
			rw.WriteHeader(http.StatusOK)
			dr := download.ToAPIRequest(downloadRequest)
			r.populateLinks(req, dr)
			encErr := encoder.Encode(dr)
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}
		} else {
			errMessage := fmt.Sprintf("Unable to find request with id:%s", requestID)
			log.Printf("server-error: %v", errMessage)

			rw.WriteHeader(http.StatusNotFound)
			encErr := encoder.Encode(errors.New(errMessage))
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}
		}
	}
}

// ValidateIncomingRequest ...
func (r *RequestResource) ValidateIncomingRequest(inReq *api.IncomingRequest) error {
	if inReq.URL == "" {
		return errors.New("empty url")
	}

	u, err := url.Parse(inReq.URL)
	if err != nil {
		return err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported url scheme: '%s'", u.Scheme)
	}
	return nil
}

// DecodeInputRequest ...
func (r *RequestResource) DecodeInputRequest(body io.Reader) (*api.IncomingRequest, error) {
	decoder := json.NewDecoder(body)
	var inReq api.IncomingRequest
	err := decoder.Decode(&inReq)
	if err != nil {
		return nil, err
	}

	return &inReq, nil
}

// GetRequestURL ...
func (r *RequestResource) GetRequestURL(id string) (*url.URL, error) {
	if r.router != nil {
		return r.router.Get("request").URL("id", id)
	}

	return nil, errors.New("no router set")
}

// Post ...
func (r *RequestResource) Post() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		apiIncomingRequest, err := r.DecodeInputRequest(req.Body)
		if err != nil {
			log.Printf("incoming-request-decode-error: %v", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = r.ValidateIncomingRequest(apiIncomingRequest)
		if err != nil {
			log.Printf("incoming-request-validation-error: %v", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		inReq := download.FromAPIIncomingRequest(apiIncomingRequest)
		downloadRequest, err := r.RequestService.ProcessNewRequest(inReq)

		if err != nil {
			log.Printf("request-processing-error: %v", err)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			encoder := json.NewEncoder(rw)
			encErr := encoder.Encode(r.WrapError(err))
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}
		} else {
			newURL, _ := r.GetRequestURL(downloadRequest.ID)
			rw.Header().Set("Content-Type", "application/json")
			rw.Header().Set("Location", newURL.String())
			rw.WriteHeader(http.StatusAccepted)
			encoder := json.NewEncoder(rw)
			dr := download.ToAPIRequest(downloadRequest)
			r.populateLinks(req, dr)
			encErr := encoder.Encode(dr)
			if encErr != nil {
				log.Printf("encode-error: %v", encErr)
			}

		}
	}
}
