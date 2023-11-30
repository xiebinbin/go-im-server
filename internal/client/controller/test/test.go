package test

type HashTwice struct {
}

type GetAuthResponse struct {
	Id           string `json:"id"`
	Token        string `json:"token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
