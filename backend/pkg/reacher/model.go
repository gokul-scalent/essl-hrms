package reacher

type VerifyEmailRequest struct {
	ToEmail string `json:"to_email"`
}

type VerifyEmailResponse struct {
	Input       string `json:"input"`
	IsReachable string `json:"is_reachable"`

	Misc struct {
		IsDisposable  bool `json:"is_disposable"`
		IsRoleAccount bool `json:"is_role_account"`
	} `json:"misc"`

	Syntax struct {
		IsValidSyntax bool `json:"is_valid_syntax"`
	} `json:"syntax"`
	//to get the provider from reacher response
	MX struct {
		AcceptsMail bool     `json:"accepts_mail"`
		Records     []string `json:"records"`
	} `json:"mx"`
}
