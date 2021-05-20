## gcal

Advanced HTTP client for golang.

## 引入http client

```
import "github.com/layatips/laya/gcal"
```

## Quick Start by trace

```
package main

import (
    "github.com/layatips/laya/gcal"
)

func main() {
	client := gcal.WithCommonHeader(genv.AppName(), "111").
		WithTrace(c, glogs.GetSpanContextKey(), url, glogs.Tracer).
		WithHeaders(map[string]string{"laya": "surprise"}).
		WithOption(gcal.OPT_TIMEOUT, 30)

	res, err := client.Get(url, map[string]string{"laya": data})
	if err != nil {
		log.Printf("BulkCourse http返回错误: err=%s", err.Error())
	}

	if res == nil {
		log.Printf("空响应")
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("BulkCourse: http状态码错误批量处理失败: %s", fmt.Sprintf("%d", res.StatusCode))
	}

	body1, err := res.ReadAll()
	log.Printf(string(body1))
}
```

## Usage

### Sending Request

```
// get
gcal.Get("http://httpbin.org/get", map[string]string{
"q": "news",
})

// get with url.Values
gcal.Get("http://httpbin.org/get", url.Values{
"q": []string{"news", "today"}
})

// post
gcal.Post("http://httpbin.org/post", map[string]string {
"name": "value"
})

// post file(multipart)
gcal.Post("http://httpbin.org/multipart", map[string]string {
"@file": "/tmp/hello.pdf",
})

// put json
gcal.PutJson("http://httpbin.org/put",
`{
    "name": "hello",
}`)

// delete
gcal.Delete("http://httpbin.org/delete")
```

### Customize Request

Before you start a new HTTP request with `Get` or `Post` method, you can specify temporary options, headers or cookies
for current request.

```
gcal.
WithHeader("User-Agent", "Super Robot").
WithHeader("custom-header", "value").
WithHeaders(map[string]string {
"another-header": "another-value",
"and-another-header": "another-value",
}).
WithOption(gcal.OPT_TIMEOUT, 60).
WithCookie(&http.Cookie{
Name: "uid",
Value: "123",
}).
Get("http://github.com")
```

### Response

The `gcal.Response` is a thin wrap of `http.Response`.

```
// traditional
res, err := gcal.Get("http://google.com")
bodyBytes, err := ioutil.ReadAll(res.Body)
res.Body.Close()

// ToString
res, err = gcal.Get("http://google.com")
bodyString, err := res.ToString()

// ReadAll
res, err = gcal.Get("http://google.com")
bodyBytes, err := res.ReadAll()
```

### Concurrent Safe

If you want to start many requests concurrently, remember to call the `Begin`
method when you begin:

```
go func () {
gcal.
Begin().
WithHeader("Req-A", "a").
Get("http://google.com")
}()
go func () {
gcal.
Begin().
WithHeader("Req-B", "b").
Get("http://google.com")
}()
```

### Error Checking

You can use `gcal.IsTimeoutError` to check for timeout error:

```
res, err := gcal.Get("http://google.com")
if gcal.IsTimeoutError(err) {
// do something
}
```

## Options

Available options as below:

- `OPT_FOLLOWLOCATION`: TRUE to follow any "Location: " header that the server sends as part of the HTTP header. Default
  to `true`.
- `OPT_CONNECTTIMEOUT`: The number of seconds or interval (with time.Duration) to wait while trying to connect. Use 0 to
  wait indefinitely.
- `OPT_CONNECTTIMEOUT_MS`: The number of milliseconds to wait while trying to connect. Use 0 to wait indefinitely.
- `OPT_MAXREDIRS`: The maximum amount of HTTP redirections to follow. Use this option alongside `OPT_FOLLOWLOCATION`.
- `OPT_PROXYTYPE`: Specify the proxy type. Valid options are `PROXY_HTTP`, `PROXY_SOCKS4`, `PROXY_SOCKS5`
  , `PROXY_SOCKS4A`. Only `PROXY_HTTP` is supported currently.
- `OPT_TIMEOUT`: The maximum number of seconds or interval (with time.Duration) to allow gcal functions to
  execute.
- `OPT_TIMEOUT_MS`: The maximum number of milliseconds to allow gcal functions to execute.
- `OPT_COOKIEJAR`: Set to `true` to enable the default cookiejar, or you can set to a `http.CookieJar` instance to use a
  customized jar. Default to `true`.
- `OPT_INTERFACE`: TODO
- `OPT_PROXY`: Proxy host and port(127.0.0.1:1080).
- `OPT_REFERER`: The `Referer` header of the request.
- `OPT_USERAGENT`: The `User-Agent` header of the request. Default to "gcal".
- `OPT_REDIRECT_POLICY`: Function to check redirect.
- `OPT_PROXY_FUNC`: Function to specify proxy.
- `OPT_UNSAFE_TLS`: Set to `true` to disable TLS certificate checking.
- `OPT_DEBUG`: Print request info.
- `OPT_CONTEXT`: Set `context.context` (can be used to cancel request).
- `OPT_BEFORE_REQUEST_FUNC`: Function to call before request is sent, option should be
  type `func(*http.Client, *http.Request)`