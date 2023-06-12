package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type RequestBatch struct {
	UUID string `json:"correlation_id"`
	URL  string `json:"original_url"`
}

type ResponseButch struct {
	Result string `json:"result"`
}
