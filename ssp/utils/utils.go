package utils

import (
	"github.com/mileusna/useragent"
	"ssp/manager"
)

func GetOsFromUA(ua string) string {
	uaObj := useragent.Parse(ua)
	return uaObj.OS
}

func GetCountryFromIP(ip string, manager manager.Manager) string {
	//ip2loc
	//return uaObj.OS
	ipClient := manager.GetIPClient()
	if ipClient != nil {
		results, err := ipClient.Get_all(ip)
		if err != nil {
			return ""
		}
		return results.Country_short
	}
	return ""
}
