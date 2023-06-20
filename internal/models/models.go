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
	UUID string `json:"correlation_id"`
	URL  string `json:"short_url"`
}

type ResponseButchForUser struct {
	SHORT_URL    string `json:"short_url"`
	ORIGINAL_URL string `json:"original_url"`
}

type Employee struct {
	Name   string
	Age    int
	Salary int
	UUID   string
}
