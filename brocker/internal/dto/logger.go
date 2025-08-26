package dto

type LogPayload struct {
	Name  string `json:"name"`
	Data  string `json:"data"`
	Level string `json:"level"`
}

type LogRPCRequestPayload struct {
	Name  string
	Data  string	
}
