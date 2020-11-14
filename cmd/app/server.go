package app

import (
	"encoding/json"
	"fmt"
	"github.com/bdaler/http/pkg/banners"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

//NewServer construct
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	log.Println("Init method")
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerById)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func (s *Server) handleGetBannerById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	log.Println("requestURI: ", request.RequestURI)

	err := request.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		log.Println(" request.ParseMultipartForm: ", err)
	}

	log.Println("request: ", request)
	log.Println("requestRequestURI: ", request.RequestURI)
	log.Println("request.FormValue(link):", request.FormValue("link"))
	//idParam := request.URL.Query().Get("id")
	//id, err := strconv.ParseInt(idParam, 10, 64)
	//if err != nil {
	//	log.Println(err)
	//	http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//	return
	//}

	image, handler, err := request.FormFile("image")
	if err != nil {
		log.Println(err)
		return
	}
	defer image.Close()
	fmt.Fprintf(writer, "%v", handler.Header)
	path, er := filepath.Abs("./web/banners/")
	if er != nil {
		log.Println(err)
	}
	i, err := os.OpenFile(path+"/1.png", os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer i.Close()
	io.Copy(i, image)

	banner := &banners.Banner{
		ID:      0,
		Title:   request.FormValue("title"),
		Content: request.FormValue("content"),
		Button:  request.FormValue("button"),
		Link:    request.FormValue("link"),
		Image:   "1.png",
	}

	item, err := s.bannersSvc.Save(request.Context(), banner)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func (s *Server) handleRemoveById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func jsonResponse(writer http.ResponseWriter, data []byte) {
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}
