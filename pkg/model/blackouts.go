package model

import (
	"github.com/kamva/mgm/v3"
	"github.com/oyamo/kplc-outage-microservice/scrapping-service/pkg/pdfutil"
)

type Blackouts struct {
	mgm.DefaultModel
	pdfutil.BlackoutResult
}
