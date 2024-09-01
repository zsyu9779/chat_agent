package endpoints

import (
	"chat_agent/config"
	"chat_agent/logger"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/zsyu9779/myUtil/validator"
	"runtime"
)

type ActionEndpoint[Req any, Res any] func(ctx context.Context, req Req) (Res, error)

// CommonBaseEndPoint endpoint 通用方法映射
func CommonBaseEndPoint[Req any, Res any](action ActionEndpoint[Req, Res]) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		typedReq, ok := req.(Req)
		if !ok {
			return nil, errors.New("invalid request type")
		}
		err = validator.Validate(typedReq)
		if err != nil {
			return nil, err
		}

		res, e := action(ctx, typedReq)
		return res, e
	}
}

// CommonBaseReflectEndPoint endpoint 类(带方法)方式支持使用
func CommonBaseReflectEndPoint[Req any, Res any](class interface {
	Method(ctx context.Context, req Req) (Res, error)
}) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		defer func() {
			isDebug := config.IsDebugEnv()
			if !isDebug {
				if err := recover(); err != nil {
					stackSlice := make([]byte, 4096)
					s := runtime.Stack(stackSlice, false)
					logger.Error(err, string(stackSlice[0:s]))
				}
			}
		}()
		typedReq, ok := req.(Req)
		if !ok {
			return nil, errors.New("invalid request type")
		}
		err = validator.Validate(typedReq)
		if err != nil {
			return nil, err
		}

		res, e := class.Method(ctx, typedReq)
		if e != nil {
			return nil, e
		}
		return res, nil
	}
}
