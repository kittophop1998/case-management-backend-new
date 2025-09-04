package utils

type CtxKey string

const (
	CtxKeyApisKey   CtxKey = "Api-Key"
	CtxKeyChannel   CtxKey = "Api-Channel"
	CtxKeyDeviceOS  CtxKey = "Api-DeviceOS"
	CtxKeyApiLang   CtxKey = "Api-Language"
	CtxKeyRequestID CtxKey = "Api-RequestID"
)
