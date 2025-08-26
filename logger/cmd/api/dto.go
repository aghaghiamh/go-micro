package main

type WriteLogRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type WriteLogRPCRequest struct {
	Name string
	Data string
}