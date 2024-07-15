package core

import "go.uber.org/zap"

func Url(url string) zap.Field {
	return zap.String("Url", url)
}

func Error(err error) zap.Field {
	return zap.NamedError("Error", err)
}

func Addr(addr string) zap.Field {
	return zap.String("Addr", addr)
}
