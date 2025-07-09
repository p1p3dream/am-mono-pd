package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
)

type Handler interface {
	TokenExchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	TokenExchangeView(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type handler struct {
	ArcDomain arc.Domain
}

type NewHandlerInput struct {
	ArcDomain arc.Domain
}

func NewHandler(in *NewHandlerInput) *handler {
	return &handler{
		ArcDomain: in.ArcDomain,
	}
}

const tokenExchangeFormHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Form Submission</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        header {
            background-color: #4a90e2;
            color: white;
            padding: 20px;
            text-align: center;
            margin-bottom: 30px;
            border-radius: 5px;
        }
        form {
            background-color: #f5f5f5;
            padding: 20px;
            border-radius: 5px;
            border: 1px solid #ddd;
        }
        .form-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        input, textarea, select {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        button {
            background-color: #4a90e2;
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover {
            background-color: #3a7bbd;
        }
    </style>
</head>
<body>
    <header>
        <h1>Data Submission Form</h1>
        <p>Please fill out the form below to submit your information</p>
    </header>

    <main>
        <form action="/zxkvkpkzhpah/token/exchange" method="post">
            <div class="form-group">
                <label for="name">API key</label>
                <input type="text" id="api-key" name="api-key" required>
            </div>

            <div class="form-group">
                <label for="name">Client Id</label>
                <input type="text" id="client-id" name="client-id" required>
            </div>

            <div class="form-group">
                <label for="name">External Id</label>
                <input type="text" id="external-id" name="external-id" required>
            </div>

            <button type="submit">Submit Form</button>
        </form>
    </main>
</body>
</html>
`

func (h *handler) TokenExchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get form values
	apiKey := r.FormValue("api-key")
	clientID := r.FormValue("client-id")
	externalID := r.FormValue("external-id")

	// Here you would typically do something with these values, like:
	// - Validating the inputs
	// - Processing the token exchange with the Arc domain
	// - Handling any errors

	if apiKey == "" || clientID == "" || externalID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	reqBody := map[string]any{
		"client_id":   clientID,
		"expire":      600,
		"external_id": externalID,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "Failed to marshal reqBody", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("POST", "https://servers-go-api.abodemine.local:23861/api/v3/auth/token/exchange", bytes.NewReader(reqBodyBytes))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send request")
		http.Error(w, "Failed to send request:"+err.Error(), http.StatusInternalServerError)
		return
	}

	if rsp.Body != nil {
		defer rsp.Body.Close()
	}

	rspBody := make(map[string]any)

	if err := json.NewDecoder(rsp.Body).Decode(&rspBody); err != nil {
		http.Error(w, "Failed to decode response body", http.StatusInternalServerError)
		return
	}

	token := rspBody["data"].(map[string]any)["token"].(string)

	http.Redirect(w, r, fmt.Sprintf("https://app.omega.test/api/auth/token/validate/%s", token), http.StatusFound)
}

func (h *handler) TokenExchangeView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(tokenExchangeFormHTML))
}
