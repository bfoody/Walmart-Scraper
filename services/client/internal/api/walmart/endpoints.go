package walmart

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/bfoody/Walmart-Scraper/utils/uuid"
)

// WalmartAPIBase is the base URL of the Walmart site
// to which endpoints are appended.
const WalmartAPIBase = "https://walmart.com"

// QuimbyAPIBase is the base URL of Walmart's quimby service,
// which provides recommendations.
const QuimbyAPIBase = "https://quimby.mobile.walmart.com"

var (
	// ItemDetailsPage returns the endpoint for viewing a single item's
	// details.
	ItemDetailsPage = func(itemSlug string, itemID string) string {
		return fmt.Sprintf("%s/ip/%s/%s", WalmartAPIBase, itemSlug, itemID)
	}
	// ItemRecommendations returns the endpoint for getting recommended items
	// for a single item.
	ItemRecommendations = func(itemID, categoryID, categoryPath, itemName string) string {
		payload := fmt.Sprintf(`{"reqId":"%s","pageId":"%s","catId":"%s","itemInfo":{"primaryCategoryPath":%s,"productName":%s,"itemAvailabilityStatus":"IN_STOCK","itemOfferType":"ONLINE_AND_STORE","isPrimaryOfferPUTEligible":true,"walledGarden":"false","verticalId":"standard"},"userReqInfo":{"referer":""},"userClientInfo":{"deviceType":"desktop","callType":"CLIENT"}}`, uuid.Generate(), itemID, categoryID, strconv.Quote(categoryPath), strconv.Quote(itemName))

		return fmt.Sprintf("%s/tempo?tenant=Walmart.com&channel=WWW&pageType=ItemPage&enrich=athenaunified,iro&item=%s&location={%%22zipCode%%22:%%2294066%%22,%%22isZipLocated%%22:false}&wm_site_mode=0&p13ncluster=&targeting={%%22nextDayStatus%%22:%%22notEligible%%22}&p13n=%s", QuimbyAPIBase, itemID, url.QueryEscape(payload))
	}
)

// url requires ?p13n=%7B%22reqId%22%3A%22920f6fff-007-17b1a44c348981%22%2C%22pageId%22%3A%22314022535%22%2C%22catId%22%3A%220%3A3944%3A1060825%3A447913%22%2C%22itemInfo%22%3A%7B%22primaryCategoryPath%22%3A%22Home%20Page%2FElectronics%2FTV%20%26%20Video%2FAll%20TVs%22%2C%22productName%22%3A%22onn.%2032%5C%22%20Class%20HD%20(720P)%20Roku%20Smart%20LED%20TV%20(100012589)%22%2C%22itemAvailabilityStatus%22%3A%22IN_STOCK%22%2C%22itemOfferType%22%3A%22ONLINE_AND_STORE%22%2C%22isPrimaryOfferPUTEligible%22%3Atrue%2C%22walledGarden%22%3A%22false%22%2C%22verticalId%22%3A%22standard%22%7D%2C%22userReqInfo%22%3A%7B%22referer%22%3A%22%22%7D%2C%22userClientInfo%22%3A%7B%22deviceType%22%3A%22desktop%22%2C%22callType%22%3A%22CLIENT%22%7D%7D
// reverse engineer ^
// {
// 	"reqId":"920f6fff-007-17b1a44c348981", // can be anything
// 	"pageId":"314022535",
// 	"catId":"0:3944:1060825:447913",
// 	"itemInfo": {
// 		"primaryCategoryPath":"Home Page/Electronics/TV & Video/All TVs",
// 		"productName":"onn. 32\" Class HD (720P) Roku Smart LED TV (100012589)",
// 		"itemAvailabilityStatus":"IN_STOCK",
// 		"itemOfferType":"ONLINE_AND_STORE",
// 		"isPrimaryOfferPUTEligible":true,
// 		"walledGarden":"false",
// 		"verticalId":"standard"
// 	},
// 	"userReqInfo": {
// 		"referer":""
// 	},
// 	"userClientInfo": {
// 		"deviceType":"desktop",
// 		"callType":"CLIENT"
// 	}
// }
