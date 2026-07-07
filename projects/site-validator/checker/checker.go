package checker

type siteChecker interface {
	isValid(url string) bool
}

type httpSiteChecker struct {
}
