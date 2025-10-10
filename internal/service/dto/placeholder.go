package dto

import "application/internal/entity"

// CreatePlaceholderReq is the request DTO for creating a placeholder.
type CreatePlaceholderReq struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// UpdatePlaceholderReq is the request DTO for updating a placeholder.
type UpdatePlaceholderReq struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// PlaceholderResp is the response DTO for a placeholder.
type PlaceholderResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToPlaceholderResp converts an entity.Placeholder to a PlaceholderResp.
func ToPlaceholderResp(e *entity.Placeholder) *PlaceholderResp {
	if e == nil {
		return nil
	}

	return &PlaceholderResp{
		ID:   e.ID.String(),
		Name: e.Name,
	}
}

type PlaceholderListResponse struct {
	Count        int                `json:"count"`
	Placeholders []*PlaceholderResp `json:"placeholders"`
}

// ToPlaceholderResps converts a slice of entity.Placeholder to a slice of PlaceholderResp.
func ToPlaceholderResps(es []entity.Placeholder) *PlaceholderListResponse {
	if es == nil {
		return nil
	}

	resps := make([]*PlaceholderResp, 0, len(es))
	for i := range es {
		resps = append(resps, ToPlaceholderResp(&es[i]))
	}

	return &PlaceholderListResponse{
		Count:        len(resps),
		Placeholders: resps,
	}
}
