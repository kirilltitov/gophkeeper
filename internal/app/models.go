package app

type requestCredentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type requestNote struct {
	Body string `json:"body" validate:"required"`
}

type requestBlob struct {
	Body string `json:"body" validate:"required"`
}

type requestTag struct {
	Tag string `json:"tag" validate:"required"`
}

type requestBankCard struct {
	Name   string `json:"name" validate:"required"`
	Number string `json:"number" validate:"required"`
	Date   string `json:"date" validate:"required"`
	CVV    string `json:"cvv" validate:"required"`
}

type baseResponse struct {
	Success bool    `json:"success"`
	Result  any     `json:"result"`
	Error   *string `json:"error"`
}
