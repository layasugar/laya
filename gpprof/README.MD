## 框架开启pprof

#### 1配置文件base栏加上pprof=true

#### 2开放端口30085

#### 3访问http://HOST:PORT/debug/pprof

#### 4火焰图本地运行go tool pprof -http=:8081 http://HOST:PORT/debug/pprof/heap

#### 5实时监控访问http://HOST:PORT/debug/charts