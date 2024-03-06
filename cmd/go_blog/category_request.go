package go_blog

type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type CategoryUpdateRequest struct {
Id int `json:"id" validate:"required"`
Name string `json:"name" validate:"required,min=1,max=100"`
}