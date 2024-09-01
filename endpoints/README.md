# endpoints
` endpoint层主要负责request／response格式的转换，以及公用拦截器相关的逻辑`
` 如：还在endpoint层提供了很多公用的拦截器，如log，metric，tracing，circuitbreaker，rate-limiter等，来保障业务系统的可用性`

**负责对transport层不同的请求协议做实现方法调度**
##如http协议请求
* endpoint 封装 handler中的请求的具体方法实现调度