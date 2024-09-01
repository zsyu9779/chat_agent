package test_endpoint

//import (
//	serverRequest "codeup.aliyun.com/aha/social-aha/aha_apiserver/service/structure/request"
//	"codeup.aliyun.com/aha/social-aha/aha_apiserver/service/test"
//	"codeup.aliyun.com/aha/social_aha_gotool/validator"
//	"context"
//	"github.com/go-kit/kit/endpoint"
//)
//
////个性化的保留用例
//
//func TestEndpoint(s test.TestService) endpoint.Endpoint {
//	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
//		r := req.(*serverRequest.Test)
//		err = validator.Validate(req)
//		if err != nil {
//			return nil, err
//		}
//		p, e := s.TestFunction(ctx, r)
//		return p, e
//	}
//}
