package main

type SendMailRequest struct {
	FromEmail string `json:"from"`
	ToEmail   string `json:"to"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}
