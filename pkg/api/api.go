package api

// BaseResponse is a response envelope containing either successful R response or an error string.
type BaseResponse[R any] struct {
	Success bool    `json:"success"` // Indicates success or failure contained in this response.
	Result  *R      `json:"result"`  // Contains a successful response payload.
	Error   *string `json:"error"`   // Contains an error message should request be a failure.
}
