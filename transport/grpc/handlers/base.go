package main

import (
	"chat_agent/config"
	"chat_agent/logger"
	"context"
	"reflect"
	"runtime"
)

type newReqFunc[GrpcReq any] func() GrpcReq
type respToGrpcResp[Res any, GrpcRes any] func(Res) GrpcRes

func CommonReflectGrpcHandler[Req any, Res any, GrpcReq any, GrpcRes any](
	class interface {
		Method(ctx context.Context, req Req) (Res, error)
	},
	newReq newReqFunc[GrpcReq],
	respConverter respToGrpcResp[Res, GrpcRes],
	validator Validator[Req],
) *transport.Server {
	errorHandler := transport.NewLogErrorHandler(logger.WriterErrorLevel())
	options := []transport.ServerOption{
		transport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			return ctx
		}),
		transport.ServerErrorHandler(errorHandler),
	}

	decodeRequest := func(ctx context.Context, grpcReq interface{}) (interface{}, error) {
		req := newReq()
		if err := convert(grpcReq.(GrpcReq), &req); err != nil {
			return nil, err
		}
		return req, nil
	}

	encodeResponse := func(_ context.Context, businessResp interface{}) (interface{}, error) {
		resp := respConverter(businessResp.(Res))
		return resp, nil
	}

	return transport.NewServer(
		CommonBaseReflectEndPoint[Req, Res](class, validator),
		decodeRequest,
		encodeResponse,
		options...,
	)
}

func convert[Src any, Dst any](source Src, target *Dst) error {
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
	if source == nil {
		return nil
	}

	sourceValue := reflect.ValueOf(source)
	targetValue := reflect.ValueOf(target).Elem()

	sourceType := sourceValue.Type()
	targetType := targetValue.Type()

	// Only convert if both are structs
	if sourceType.Kind() == reflect.Struct && targetType.Kind() == reflect.Struct {
		for i := 0; i < sourceType.NumField(); i++ {
			sourceField := sourceType.Field(i)
			if !sourceField.IsExported() {
				continue
			}
			sourceFieldValue := sourceValue.Field(i)

			if targetField := targetValue.FieldByName(sourceField.Name); targetField.IsValid() && targetField.CanSet() {
				targetField.Set(sourceFieldValue)
			}
		}
	}
	return nil
}
