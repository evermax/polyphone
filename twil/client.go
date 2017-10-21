package twil

// TwilioClient
type TwilioClient struct {
	accountSid    string
	accountSecret string
	appSid        string
	name          string
	number        string
}

// NewTwilioClient creates a TwilioClient ussing the parameters
func NewTwilioClient(accountSid, accountSecret, appSid, name, number string) *TwilioClient {
	return &TwilioClient{
		accountSid:    accountSid,
		accountSecret: accountSecret,
		appSid:        appSid,
		name:          name,
		number:        number,
	}
}
