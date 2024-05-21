package httprequest

type HeaderKey string

const (
	RequestId      HeaderKey = "X-Request-ID"
	ContentTypeKey HeaderKey = "Content-Type"
	Authorization  HeaderKey = "Authorization"
)

type ContentType string

const (
	ApplicationJson ContentType = "application/json"
)

type Method string

const (
	Options Method = "OPTIONS"
	Get     Method = "GET"
	Head    Method = "HEAD"
	Post    Method = "POST"
	Put     Method = "PUT"
	Delete  Method = "DELETE"
	Trace   Method = "TRACE"
	Connect Method = "CONNECT"
)
