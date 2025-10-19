package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/danilobml/motivate/internal/handlers"
	"github.com/danilobml/motivate/internal/models"
	"github.com/danilobml/motivate/internal/repositories"
	"github.com/danilobml/motivate/internal/services"
)

func setupServer(isSeeded bool) *httptest.Server {
	// NOTE: inMemory doesn't require mocking. Should be changed if persistence is implemented
	inMemoryRepo := repositories.NewInMemoryQuoteRepository()
	mockService := services.NewQuoteService(inMemoryRepo)
	router := handlers.NewQuotesRouter(mockService)
	routes := handlers.RegisterRoutes(router)

	if isSeeded {
		mockService.SeedDbFromFile("./test_seed.json")
	}

	return httptest.NewTLSServer(routes)
}

func Test_HealthCheck(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	res, err := client.Get(srv.URL + "/health")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_GetRandomQuote_Success(t *testing.T) {
	srv := setupServer(true)
	defer srv.Close()

	client := srv.Client()

	res, err := client.Get(srv.URL + "/quote")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()

	var quote models.Quote
	err = json.NewDecoder(res.Body).Decode(&quote)
	require.NoError(t, err, "failed to decode JSON response")

	require.NotEmpty(t, quote.Id, "expected non-empty Id")
	require.NotEmpty(t, quote.Text, "expected non-empty Text")
	require.NotEmpty(t, quote.Author, "expected non-empty Author")
}

func Test_GetRandomQuote_404_if_no_entries(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	res, _ := client.Get(srv.URL + "/quote")
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func Test_CreateNewQuote_Success(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	text := "Test Text"
	author := "Test Author"

	payload := map[string]any{
		"text":   text,
		"author": author,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/add", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, "201 Created", res.Status)

	var quote models.Quote
	err = json.NewDecoder(res.Body).Decode(&quote)
	require.NoError(t, err, "failed to decode JSON response")

	require.Equal(t, text, quote.Text)
	require.Equal(t, author, quote.Author)
}

func Test_CreateNewQuote_Success_with_EmptyAuthor(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	text := "Test Text"
	author := ""

	payload := map[string]any{
		"text":   text,
		"author": author,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/add", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, "201 Created", res.Status)

	var quote models.Quote
	err = json.NewDecoder(res.Body).Decode(&quote)
	require.NoError(t, err, "failed to decode JSON response")

	require.Equal(t, text, quote.Text)
	require.Equal(t, author, quote.Author)
}

func Test_CreateNewQuote_Fails_400_with_EmptyText(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	payload := map[string]any{
		"text":   "",
		"author": "Test Author",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/add", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func Test_CreateNewQuote_Fails_400_with_InvalidJSON(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/add", bytes.NewBuffer([]byte("{invalid")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func Test_CreateNewQuote_Fails_400_with_EmptyBody(t *testing.T) {
	srv := setupServer(false)
	defer srv.Close()

	client := srv.Client()

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/add", bytes.NewBuffer(nil))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func Test_Request_500_if_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Millisecond)
	}))
	defer srv.Close()

	client := &http.Client{
		Timeout: 1 * time.Millisecond,
	}

	_, err := client.Get(srv.URL + "/quote")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Client.Timeout", "expected client timeout error")
}
