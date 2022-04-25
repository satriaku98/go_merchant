package apprequest

type CustomerUpdateRequest struct {
	Name string `json:"name" bind:"required"`
}
