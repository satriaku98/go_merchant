package apprequest

type CustomerRequest struct {
	Name     string `json:"name" bind:"required"`
	Passcode string `json:"passcode" bind:"required"`
}
