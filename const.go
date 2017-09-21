package lapi

const (
	// Error code format: x.xxx.xxx?x
	// First 0 is used by system.
	// Next 001 - 100 is used by lapi.
	// Other sub-package could obtain range 100 - 999
	// Last 3 digits for error code, it might be extended to
	// whatever number if necessary

	// App errors
	ERR_ROUTER_NOT_DEFINED    = "0.001.001"
	ERR_SERVER_CONFIG_MISSING = "0.001.002"
	ERR_NO_HANDLER_FOUND      = "0.001.003"

	// Router, http errors
	ERR_HTTP_NOT_FOUND   = "0.002.001"
	ERR_HTTP_BAD_REQUEST = "0.002.002"

	// Request, Response, Body, Parser errors
	ERR_RESPONSE_ALREADY_SENT = "0.003.001"
	ERR_NO_PARSER_FOUND       = "0.003.002"
	ERR_CONTENT_TYPE_EMPTY    = "0.003.003"
	ERR_NO_WRITER_FOUND       = "0.003.004"
	ERR_PARSE_INVALID_CONTENT = "0.003.005"

	// Container error
	ERR_BIND_INVALID_INTERFACE         = "0.004.001"
	ERR_BIND_INVALID_CONCRETE          = "0.004.002"
	ERR_BIND_NOT_IMPLEMENT_INTERFACE   = "0.004.003"
	ERR_BIND_INVALID_STRUCT            = "0.004.004"
	ERR_BIND_INVALID_STRUCT_CONCRETE   = "0.004.005"
	ERR_BIND_INVALID_ARGUMENTS         = "0.004.006"
	ERR_RESOLVE_NOT_EXIST_ABSTRACT     = "0.004.007"
	ERR_RESOLVE_INVALID_CONCRETE       = "0.004.008"
	ERR_RESOLVE_INSUFFICIENT_ARGUMENTS = "0.004.009"
	ERR_RESOLVE_NON_VALUES_RETURNED    = "0.004.010"
	ERR_RESOLVE_INVALID_ARGUMENTS      = "0.004.011"
	ERR_INJECT_INVALID_TARGET_TYPE     = "0.004.012"

	PORT_HTTP  = 80
	PORT_HTTPS = 443

	SCHEME_HTTP  = "http"
	SCHEME_HTTPS = "https"

	HEADER_CONTENT_TYPE = "content-type"

	CONTENT_TYPE_JSON       = "application/json"
	CONTENT_TYPE_XML        = "application/xml"
	CONTENT_TYPE_TEXT       = "text/plain"
	CONTENT_TYPE_DEFAULT    = CONTENT_TYPE_JSON
	CONTENT_CHARSET_DEFAULT = "utf-8"
)
