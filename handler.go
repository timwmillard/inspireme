package inspireme

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// Handler -
type Handler struct {
	Log       *log.Logger
	InspireMe *ImageGenerator
}

// NewHandler -
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP handles the routing
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// GET method
	switch r.Method {
	case http.MethodPost:
		h.generate(w, r)
	}
}

type postRequest struct {
	Quote         string `json:"quote"`
	BackgroundURL string `json:"backgroundUrl"`
}

type postResponce struct {
	ID       string `json:"id,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

// Gererate an image
// quote and backgroundUrl
func (h *Handler) generate(w http.ResponseWriter, r *http.Request) {

	var req postRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&req)
	if err != nil {
		h.Log.Printf("unable to decode json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.Printf("generate image \"%v\" %v", req.Quote, req.BackgroundURL)

	ctx := context.TODO()
	imgURL, err := h.InspireMe.GenerateAndStore(ctx, req.Quote, req.BackgroundURL, nil)
	if err != nil {
		h.Log.Printf("unable to generate image: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// imgURL := "https://miro.medium.com/max/3152/1*Ifpd_HtDiK9u6h68SZgNuA.png"

	var resp postResponce
	resp.ImageURL = imgURL
	encoder := json.NewEncoder(w)
	encoder.Encode(&resp)
}
