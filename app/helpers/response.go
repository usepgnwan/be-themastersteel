package helpers

import "os"

type Response struct {
	Message string      `json:"message" form:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
}
type ResponsePaginate struct {
	Total       int64       `json:"total"`
	Rows        interface{} `json:"rows"`
	CurrentPage int         `json:"currentPage"`
	PerPage     int         `json:"perPage"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	LastPage    int         `json:"lastPage"`
}

var AppsGroup = []string{"layer1&2", "layer3"}
var SecretkeySeller = os.Getenv("POOLAPACK_HEADER_KEY")
