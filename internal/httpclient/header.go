package httpclient

const (
	HeaderContentType   = "Content-Type"
	HeaderAccept        = "Accept"
	HeaderAuthorization = "Authorization"
	HeaderAPIKey        = "X-API-Key"

	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

func JSONHeaders() map[string]string {
	return map[string]string{
		HeaderContentType: ContentTypeJSON,
		HeaderAccept:      ContentTypeJSON,
	}
}

func FormURLEncodedHeaders() map[string]string {
	return map[string]string{
		HeaderContentType: ContentTypeFormURLEncoded,
	}
}
