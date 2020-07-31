package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	config2 "wooc/config"
	"wooc/middlewares"
	"wooc/models"
)

var srv *httptest.Server

func TestMain(m *testing.M) {
	r := mux.NewRouter()

	config := config2.Config{
		PageSize: 2,
		Tokens:   []string{"0123456789"},
	}
	ipHandler := NewIpHandler(config)
	r.HandleFunc("/ip", ipHandler.GetList).Methods(http.MethodGet)
	r.HandleFunc("/ip", ipHandler.Add).Methods(http.MethodPost)

	amw := middlewares.AuthMiddleware{Tokens: config.Tokens}
	r.Use(amw.Middleware)

	srv = httptest.NewServer(r)
	defer srv.Close()

	os.Exit(m.Run())
}

func TestIpHandler_Add(t *testing.T) {
	jsonStr := []byte(`{
 		"ip_address": "8.8.8.8",
 		"ASN": 1234,
		"domains": ["alpha.com", "beta.com"]
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ip", srv.URL), bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 0123456789")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var ipModel models.IpModel
	if err = json.Unmarshal(resBytes, &ipModel); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1234,
		Domains:   []string{"alpha.com", "beta.com"},
	}, ipModel)
}

func TestIpHandler_AddNotUnauthorized(t *testing.T) {
	jsonStr := []byte(`{
 		"ip_address": "8.8.8.8",
 		"ASN": 1234,
		"domains": ["alpha.com", "beta.com"]
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/ip", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestIpHandler_AddBadRequest(t *testing.T) {
	jsonStr := []byte(`{
 		"ip_address": "8.8.8.8",
 		"ASN": "1234",
		"domains": ["alpha.com", "beta.com"]
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ip", srv.URL), bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 0123456789")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestIpHandler_AddNotValid(t *testing.T) {
	jsonStr := []byte(`{
 		"ip_address": "8.8.8.8",
 		"ASN": 1234,
		"domains": ["alpha", "beta.com"]
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ip", srv.URL), bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 0123456789")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestIpHandler_GetList(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ip?page=1", srv.URL), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer 0123456789")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var ipModels []*models.IpModel
	if err = json.Unmarshal(resBytes, &ipModels); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1234,
		Domains:   []string{"alpha.com", "beta.com"},
	}, ipModels[0])
}

func TestIpHandler_GetListNotUnauthorized(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("%s/ip?page=1", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var ipModels []*models.UnauthorizedIpModel
	if err = json.Unmarshal(resBytes, &ipModels); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &models.UnauthorizedIpModel{IpAddress: "8.8.8.8"}, ipModels[0])
}

func TestIpHandler_GetListWithoutPageParam(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ip", srv.URL), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer 0123456789")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var ipModels []*models.IpModel
	if err = json.Unmarshal(resBytes, &ipModels); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1234,
		Domains:   []string{"alpha.com", "beta.com"},
	}, ipModels[0])
}

func TestIpHandler_GetListBadRequest(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("%s/ip?page", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestIpHandler_GetListZeroPage(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("%s/ip?page=0", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}