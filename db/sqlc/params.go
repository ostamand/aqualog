package db

import "time"

const defaultParamsLimit = 50

func (params *GetParamsByTypeParams) FillDefaults() {
	if params.To.IsZero() {
		// taking the guess that 100 years added to current time should be enough
		params.To = time.Now().Add(time.Hour * 24 * 365 * 100)
	}
	if params.Limit == 0 {
		params.Limit = defaultParamsLimit
	}
}
