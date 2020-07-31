package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"wooc/config"
	"wooc/interfaces"
	"wooc/models"
	"wooc/repositories"
)

type IpHandler struct {
	validator  *validator.Validate
	repository interfaces.RepositoryInterface
	config     config.Config
}

func NewIpHandler(config config.Config) *IpHandler {
	return &IpHandler{
		validator:  validator.New(),
		repository: repositories.NewIpMemoryRepository(),
		config:     config,
	}
}

func (h *IpHandler) GetList(w http.ResponseWriter, r *http.Request) {
	pages, ok := r.URL.Query()["page"]
	if !ok || len(pages) < 1 {
		pages = []string{"1"}
	}
	page, err := strconv.ParseUint(pages[0], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ipModels, err := h.repository.FindAll(page, h.config.PageSize)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Header.Get("Authorization") != "unauthorized" {
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipModels); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var unauthorizedIpModels []models.UnauthorizedIpModel
	for k := range ipModels {
		unauthorizedIpModels = append(
			unauthorizedIpModels,
			models.UnauthorizedIpModel{IpAddress: ipModels[k].IpAddress},
		)
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(unauthorizedIpModels); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *IpHandler) Add(w http.ResponseWriter, r *http.Request) {
	var ipModel models.IpModel
	if err := json.NewDecoder(r.Body).Decode(&ipModel); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(&ipModel); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newIp := h.repository.Create(ipModel)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newIp); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
