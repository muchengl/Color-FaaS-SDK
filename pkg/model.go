package cfaas

type FuncRequest string
type FuncResponse string

type funcInvokeRequest struct {
	Msg string `json:"umsg"`
}

type funcInvokeResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
