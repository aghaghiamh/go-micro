package domain

type EmailServer struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	ToEmail     string
	FromEmail   string
	FromName    string
	Subject     string
	Attachments []string // path to each attachment
	Data        any
	DataMap     map[string]any
}
