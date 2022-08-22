package utils

func BuildEndpointUrl(url string, port int) string {

	var baseUrl string

	if url == "" {
		baseUrl = "https://upload.bugsnag.com"
	} else {
		baseUrl = url
	}

	if port != 0 {
		return baseUrl + ":" + string(port)
	}

	return baseUrl
}
