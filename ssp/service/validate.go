package service

import (
	"context"
	"errors"
	"ssp/model"
)

func (s *SSPSvcHandler) ValidateGetAdMarkupPayload(ctx context.Context, request *model.AdRequest) error {
	if request.AdSlotId == "" {
		return errors.New("adslotid  missing in GetAdMarkupPayload")
	} else if request.AdType == "" {
		return errors.New("adtype missing in GetAdMarkupPayload")
	}
	return nil
}
