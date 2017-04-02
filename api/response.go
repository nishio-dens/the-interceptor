package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	mimeXml = "application/xml"
)

func SendSuccess(w http.ResponseWriter, response []byte, mime string) {
	writeCommonHeaders(w)
	w.Header().Set("Content-Type", string(mime))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func SendNoBodySuccess(w http.ResponseWriter) {
	writeCommonHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func SendInternalError(w http.ResponseWriter, response []byte, mime string) {
	writeCommonHeaders(w)
	w.Header().Set("Content-Type", string(mime))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(response)
}

func SendBadRequestError(w http.ResponseWriter, response []byte, mime string) {
	writeCommonHeaders(w)
	w.Header().Set("Content-Type", string(mime))
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}

func SendNotFoundError(w http.ResponseWriter, response []byte, mime string) {
	writeCommonHeaders(w)
	w.Header().Set("Content-Type", string(mime))
	w.WriteHeader(http.StatusNotFound)
	w.Write(response)
}

func SendSuccessXml(w http.ResponseWriter, response interface{}) {
	r, err := xml.MarshalIndent(response, "", "    ")
	if err != nil {
		SendInternalError(w, []byte("Internal Server Error"), string(mimeXml))
		return
	}
	w.Header().Set("Content-Type", string(mimeXml))
	SendSuccess(w, []byte(r), mimeXml)
}

func SendInternalErrorXml(w http.ResponseWriter, response interface{}) {
	r, err := xml.MarshalIndent(response, "", "    ")
	if err != nil {
		SendInternalError(w, []byte("Internal Server Error"), string(mimeXml))
		return
	}
	SendInternalError(w, []byte(r), string(mimeXml))
}

func SendBadRequestXml(w http.ResponseWriter, response interface{}) {
	r, err := xml.MarshalIndent(response, "", "    ")
	if err != nil {
		SendInternalError(w, []byte("Internal Server Error"), string(mimeXml))
		return
	}
	w.Header().Set("Content-Type", string(mimeXml))
	SendBadRequestError(w, []byte(r), string(mimeXml))
}

func SendNotFoundXml(w http.ResponseWriter, response interface{}) {
	r, err := xml.MarshalIndent(response, "", "    ")
	if err != nil {
		SendInternalError(w, []byte("Internal Server Error"), string(mimeXml))
		return
	}
	w.Header().Set("Content-Type", string(mimeXml))
	SendNotFoundError(w, []byte(r), string(mimeXml))
}

// Private

func writeCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("X-The-Interceptor", fmt.Sprintf("Version %s", CurrentApiVersion))
}
