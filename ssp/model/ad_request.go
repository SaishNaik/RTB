package model

import "errors"

//todo update this struct with additional properties
type AdRequest struct {
	AdSlotId string `json:"adslotid"`
	AdType   string `json:"adtype"`
	PubId    string `json:"pubid"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	UA       string `json:"ua"`
	IP       string `json:"ip"`
}

func (ar *AdRequest) Validate() error {
	if ar.AdType == "" {
		return errors.New("ad type not sent in AdRequest")
	} else if ar.AdSlotId == "" {
		return errors.New("ad slot id not sent in AdRequest")
	}
	return nil
}
