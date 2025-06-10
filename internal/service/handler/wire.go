package handler

import (
	"application/internal/service"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(

	NewServiceList,
	NewMuxHealthzHandler,

	NewMuxTokenHandler,
	NewMuxCampaignHandler,
)

// New ServiceList
func NewServiceList(
	healthzSvc *HealthzHandler,

	tokenSvc *TokenHandler,
	campaignSvc *CampaignHandler,
) []service.Handler {
	return []service.Handler{
		healthzSvc,
		tokenSvc,
		campaignSvc,
	}
}
