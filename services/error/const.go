package error

// error code
const (
	ResultSuccess    = "000000" //请求成功
	ParamsErrorCode  = "000001" //参数错误
	TokenExpired     = "000002" //jwt token 过期
	TokenInvalid     = "000003" //jwt token 不可用
	MidddleErrorCode = "000004" //中间件错误

	ServerErrorCode = "500000" //mysql,redis,网络通信的错误

	GrpcResourceOssGetStsError = "11000001" // 远程调用oss getsts 错误

)
