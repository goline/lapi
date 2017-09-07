package lapi

const (
	ERROR_ROUTER_NOT_DEFINED    = 1
	ERROR_SERVER_CONFIG_MISSING = 2
	ERROR_HTTP_NOT_FOUND        = 3
	ERROR_RESPONSE_ALREADY_SENT = 4
	ERROR_NO_PARSER_FOUND       = 5
	ERROR_HTTP_BAD_REQUEST      = 6
	ERROR_NO_HANDLER_FOUND      = 7

	PORT_HTTP  = 80
	PORT_HTTPS = 443

	SCHEME_HTTP  = "http"
	SCHEME_HTTPS = "https"

	CONTENT_TYPE_JSON    = "application/json"
	CONTENT_TYPE_XML     = "application/xml"
	CONTENT_TYPE_DEFAULT = CONTENT_TYPE_JSON
)
