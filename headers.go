package inpu

const (
	// Content and Accept headers
	HeaderAccept          = "Accept"
	HeaderAcceptCharset   = "Accept-Charset"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderAcceptLanguage  = "Accept-Language"
	HeaderContentType     = "Content-Type"
	HeaderContentLength   = "Content-Length"
	HeaderContentEncoding = "Content-Encoding"

	// Authentication headers
	HeaderAuthorization = "Authorization"
	HeaderAPIKey        = "X-API-Key"
	HeaderAPISecret     = "X-API-Secret"
	HeaderAPIToken      = "X-API-Token"

	// Client identification
	HeaderUserAgent = "User-Agent"
	HeaderReferer   = "Referer"
	HeaderOrigin    = "Origin"

	// Caching and conditional requests
	HeaderCacheControl      = "Cache-Control"
	HeaderIfModifiedSince   = "If-Modified-Since"
	HeaderIfNoneMatch       = "If-None-Match"
	HeaderIfMatch           = "If-Match"
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	HeaderIfRange           = "If-Range"

	// Request control
	HeaderExpect     = "Expect"
	HeaderRange      = "Range"
	HeaderHost       = "Host"
	HeaderConnection = "Connection"
	HeaderUpgrade    = "Upgrade"
	HeaderTE         = "TE"

	// Cookies and session
	HeaderCookie = "Cookie"

	// Request tracking and tracing
	HeaderXRequestID     = "X-Request-ID"
	HeaderXCorrelationID = "X-Correlation-ID"
	HeaderXTraceID       = "X-Trace-ID"
	HeaderXSpanID        = "X-Span-ID"

	// Custom client headers
	HeaderXClientVersion = "X-Client-Version"
	HeaderXClientName    = "X-Client-Name"
	HeaderXForwardedFor  = "X-Forwarded-For"
	HeaderXRealIP        = "X-Real-IP"

	// Security headers (client-side)
	HeaderXCSRFToken     = "X-CSRF-Token"
	HeaderXRequestedWith = "X-Requested-With"

	// CORS preflight headers
	HeaderAccessControlRequestMethod  = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"

	// Common API headers
	HeaderXAPIVersion = "X-API-Version"
	HeaderXClientID   = "X-Client-ID"
	HeaderXSessionID  = "X-Session-ID"
	HeaderXTenantID   = "X-Tenant-ID"
)

const (
	// Text types
	MimeTypeText       = "text/plain"
	MimeTypeHtml       = "text/html"
	MimeTypeCss        = "text/css"
	MimeTypeJavascript = "text/javascript"
	MimeTypeCsv        = "text/csv"
	MimeTypeTextXml    = "text/xml"
	MimeTypeCalendar   = "text/calendar"

	// Application types
	MimeTypeJson              = "application/json"
	MimeTypeApplicationXml    = "application/xml"
	MimeTypeFormUrlEncoded    = "application/x-www-form-urlencoded"
	MimeTypeMultipartFormData = "multipart/form-body"
	MimeTypeOctetStream       = "application/octet-stream"
	MimeTypePdf               = "application/pdf"
	MimeTypeZip               = "application/zip"
	MimeTypeGzip              = "application/gzip"
	MimeTypeTar               = "application/x-tar"
	MimeTypeRar               = "application/vnd.rar"
	MimeType7z                = "application/x-7z-compressed"
	MimeTypeJsonApi           = "application/vnd.api+json"
	MimeTypeJsonPatch         = "application/json-patch+json"
	MimeTypeJsonMergePatch    = "application/merge-patch+json"

	// Image types
	MimeTypeJpeg = "image/jpeg"
	MimeTypePng  = "image/png"
	MimeTypeGif  = "image/gif"
	MimeTypeWebp = "image/webp"
	MimeTypeSvg  = "image/svg+xml"
	MimeTypeBmp  = "image/bmp"
	MimeTypeIco  = "image/x-icon"
	MimeTypeTiff = "image/tiff"
	MimeTypeAvif = "image/avif"

	// Audio types
	MimeTypeMp3  = "audio/mpeg"
	MimeTypeWav  = "audio/wav"
	MimeTypeOgg  = "audio/ogg"
	MimeTypeAac  = "audio/aac"
	MimeTypeFlac = "audio/flac"
	MimeTypeM4a  = "audio/mp4"

	// Video types
	MimeTypeMp4  = "video/mp4"
	MimeTypeAvi  = "video/x-msvideo"
	MimeTypeMov  = "video/quicktime"
	MimeTypeWmv  = "video/x-ms-wmv"
	MimeTypeFlv  = "video/x-flv"
	MimeTypeWebm = "video/webm"
	MimeTypeMkv  = "video/x-matroska"

	// Font types
	MimeTypeWoff  = "font/woff"
	MimeTypeWoff2 = "font/woff2"
	MimeTypeTtf   = "font/ttf"
	MimeTypeOtf   = "font/otf"
	MimeTypeEot   = "application/vnd.ms-fontobject"

	// Microsoft Office
	MimeTypeDocx = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MimeTypeDoc  = "application/msword"
	MimeTypeXlsx = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MimeTypeXls  = "application/vnd.ms-excel"
	MimeTypePptx = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MimeTypePpt  = "application/vnd.ms-powerpoint"

	// Other common types
	MimeTypeRtf  = "application/rtf"
	MimeTypeSwf  = "application/x-shockwave-flash"
	MimeTypeAtom = "application/atom+xml"
	MimeTypeRss  = "application/rss+xml"
	MimeTypeYaml = "application/x-yaml"
	MimeTypeToml = "application/toml"
)
