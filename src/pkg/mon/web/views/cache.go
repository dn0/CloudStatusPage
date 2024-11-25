package views

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	cacheHeader        = "Cache-Control"
	cachePrivate       = "private"
	cachePublic        = "public"
	CacheMaxAgeDefault = 30 * time.Second
)

//nolint:gochecknoglobals // These are constants.
var (
	CachingDisabled       = false // I know, this is ugly, but it's set only once during startup.
	CacheAllowedAnyValues = CacheAllowedQueryValues{}
)

type CacheAllowedQueryValues = map[string]struct{}

type CacheAllowedQueryParams = map[string]CacheAllowedQueryValues

func CacheControl(w http.ResponseWriter, r *http.Request, maxAge time.Duration, whitelist CacheAllowedQueryParams) {
	if CachingDisabled {
		return
	}

	directive := cachePrivate
	if isStaticPage(r.URL.Query(), whitelist) {
		directive = cachePublic
	}
	directive += ", max-age=" + strconv.Itoa(int(maxAge.Seconds()))

	w.Header().Set(cacheHeader, directive)
}

// isStaticPage determines whether a page can be cached by a proxy/CDN.
// Definition of a static page in this project:
// - _No_ query parameters: /some/url,
// - _Only_ allowed query parameters with arbitrary values: /some/url?allowed_param=anything_goes,
// - _Only_ allowed query parameters with allowed values: /some/url?allowed_param=allowed_value.
// Everything else is considered a dynamic page and should not be cached at all or can be cached by a browser.
// Example whitelist:
//
//	CacheAllowedQueryParams{
//		"page": CacheAllowedAnyValues, // any value is OK => the view is responsible for validating it;
//	                                                      (a view error should remove the Cache-Control header)
//		"status": CacheAllowedQueryValues{
//			"0": struct{}{},
//			"1": struct{}{},
//			"2": struct{}{},
//	 },
//	}
func isStaticPage(qs url.Values, whitelist CacheAllowedQueryParams) bool {
	if whitelist == nil {
		return len(qs) == 0
	}
	for qsKey, qsValues := range qs {
		allowedValues, found := whitelist[qsKey]
		if !found {
			return false // query string parameter now allowed
		}
		if len(allowedValues) > 0 {
			for _, v := range qsValues {
				if _, ok := allowedValues[v]; !ok {
					return false // query string parameter value not allowed
				}
			}
		}
	}
	return true
}
