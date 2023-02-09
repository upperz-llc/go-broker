package domain

type SuperuserCheckPOST struct {
	Username string `json:"username"`
}

type ClientCheckPOST struct {
	ClientID string `json:"clientid"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type ACLCheckPOST struct {
	Username string `json:"username"`
	ClientID string `json:"clientid"`
	Topic    string `json:"topic"`
	ACC      string `json:"acc"`
}
