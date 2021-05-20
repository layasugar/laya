// Powerful and easy to use http client
package gcal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gutils"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"log"
	"strings"

	"time"

	"io"
	"io/ioutil"
	"sync"

	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"

	"crypto/tls"

	"compress/gzip"

	"encoding/json"
	"mime/multipart"
)

// Constants definations
// CURL options, see https://github.com/bagder/curl/blob/169fedbdce93ecf14befb6e0e1ce6a2d480252a3/packages/OS400/curl.inc.in
const (
	USERAGENT = "gcal"
)

const (
	PROXY_HTTP int = iota
	PROXY_SOCKS4
	PROXY_SOCKS5
	PROXY_SOCKS4A

	// CURL like OPT
	OPT_AUTOREFERER
	OPT_FOLLOWLOCATION
	OPT_CONNECTTIMEOUT
	OPT_CONNECTTIMEOUT_MS
	OPT_MAXREDIRS
	OPT_PROXYTYPE
	OPT_TIMEOUT
	OPT_TIMEOUT_MS
	OPT_COOKIEJAR
	OPT_INTERFACE
	OPT_PROXY
	OPT_REFERER
	OPT_USERAGENT

	// Other OPT
	OPT_REDIRECT_POLICY
	OPT_PROXY_FUNC
	OPT_DEBUG
	OPT_UNSAFE_TLS

	OPT_CONTEXT

	OPT_BEFORE_REQUEST_FUNC
	OPT_AFTER_REQUEST_FUNC
)

// String map of options
var CONST = map[string]int{
	"OPT_AUTOREFERER":         OPT_AUTOREFERER,
	"OPT_FOLLOWLOCATION":      OPT_FOLLOWLOCATION,
	"OPT_CONNECTTIMEOUT":      OPT_CONNECTTIMEOUT,
	"OPT_CONNECTTIMEOUT_MS":   OPT_CONNECTTIMEOUT_MS,
	"OPT_MAXREDIRS":           OPT_MAXREDIRS,
	"OPT_PROXYTYPE":           OPT_PROXYTYPE,
	"OPT_TIMEOUT":             OPT_TIMEOUT,
	"OPT_TIMEOUT_MS":          OPT_TIMEOUT_MS,
	"OPT_COOKIEJAR":           OPT_COOKIEJAR,
	"OPT_INTERFACE":           OPT_INTERFACE,
	"OPT_PROXY":               OPT_PROXY,
	"OPT_REFERER":             OPT_REFERER,
	"OPT_USERAGENT":           OPT_USERAGENT,
	"OPT_REDIRECT_POLICY":     OPT_REDIRECT_POLICY,
	"OPT_PROXY_FUNC":          OPT_PROXY_FUNC,
	"OPT_DEBUG":               OPT_DEBUG,
	"OPT_UNSAFE_TLS":          OPT_UNSAFE_TLS,
	"OPT_CONTEXT":             OPT_CONTEXT,
	"OPT_BEFORE_REQUEST_FUNC": OPT_BEFORE_REQUEST_FUNC,
	"OPT_AFTER_REQUEST_FUNC":  OPT_AFTER_REQUEST_FUNC,
}

// Default options for any clients.
var defaultOptions = map[int]interface{}{
	OPT_FOLLOWLOCATION: true,
	OPT_MAXREDIRS:      10,
	OPT_AUTOREFERER:    true,
	OPT_USERAGENT:      USERAGENT,
	OPT_COOKIEJAR:      true,
	OPT_DEBUG:          false,
}

// These options affect transport, transport may not be reused if you change any
// of these options during a request.
var transportOptions = []int{
	OPT_CONNECTTIMEOUT,
	OPT_CONNECTTIMEOUT_MS,
	OPT_PROXYTYPE,
	OPT_TIMEOUT,
	OPT_TIMEOUT_MS,
	OPT_INTERFACE,
	OPT_PROXY,
	OPT_PROXY_FUNC,
	OPT_UNSAFE_TLS,
}

// These options affect cookie jar, jar may not be reused if you change any of
// these options during a request.
var jarOptions = []int{
	OPT_COOKIEJAR,
}

// Thin wrapper of http.Response(can also be used as http.Response).
type Response struct {
	*http.Response
}

// Read response body into a byte slice.
func (client *Response) ReadAll() ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch client.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(client.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = client.Body
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// Read response body into string.
func (client *Response) ToString() (string, error) {
	bytes, err := client.ReadAll()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Prepare a request.
func prepareRequest(method string, url_ string, headers map[string]string,
	body io.Reader, options map[int]interface{}) (*http.Request, error) {
	req, err := http.NewRequest(method, url_, body)

	if err != nil {
		return nil, err
	}

	// OPT_REFERER
	if referer, ok := options[OPT_REFERER]; ok {
		if refererStr, ok := referer.(string); ok {
			req.Header.Set("Referer", refererStr)
		}
	}

	// OPT_USERAGENT
	if useragent, ok := options[OPT_USERAGENT]; ok {
		if useragentStr, ok := useragent.(string); ok {
			req.Header.Set("User-Agent", useragentStr)
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Prepare a transport.
//
// Handles timemout, proxy and maybe other transport related options here.
func prepareTransport(options map[int]interface{}) (http.RoundTripper, error) {
	transport := &http.Transport{}

	var connectTimeout time.Duration

	if connectTimeoutMS_, ok := options[OPT_CONNECTTIMEOUT_MS]; ok {
		if connectTimeoutMS, ok := connectTimeoutMS_.(int); ok {
			connectTimeout = time.Duration(connectTimeoutMS) * time.Millisecond
		} else {
			return nil, fmt.Errorf("OPT_CONNECTTIMEOUT_MS must be int")
		}
	} else if connectTimeout_, ok := options[OPT_CONNECTTIMEOUT]; ok {
		if connectTimeout, ok = connectTimeout_.(time.Duration); !ok {
			if connectTimeoutS, ok := connectTimeout_.(int); ok {
				connectTimeout = time.Duration(connectTimeoutS) * time.Second
			} else {
				return nil, fmt.Errorf("OPT_CONNECTTIMEOUT must be int or time.Duration")
			}
		}
	}

	var timeout time.Duration

	if timeoutMS_, ok := options[OPT_TIMEOUT_MS]; ok {
		if timeoutMS, ok := timeoutMS_.(int); ok {
			timeout = time.Duration(timeoutMS) * time.Millisecond
		} else {
			return nil, fmt.Errorf("OPT_TIMEOUT_MS must be int")
		}
	} else if timeout_, ok := options[OPT_TIMEOUT]; ok {
		if timeout, ok = timeout_.(time.Duration); !ok {
			if timeoutS, ok := timeout_.(int); ok {
				timeout = time.Duration(timeoutS) * time.Second
			} else {
				return nil, fmt.Errorf("OPT_TIMEOUT must be int or time.Duration")
			}
		}
	}

	// fix connect timeout(important, or it might cause a long time wait during
	//connection)
	if timeout > 0 && (connectTimeout > timeout || connectTimeout == 0) {
		connectTimeout = timeout
	}

	transport.Dial = func(network, addr string) (net.Conn, error) {
		var conn net.Conn
		var err error
		if connectTimeout > 0 {
			conn, err = net.DialTimeout(network, addr, connectTimeout)
			if err != nil {
				return nil, err
			}
		} else {
			conn, err = net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
		}

		if timeout > 0 {
			conn.SetDeadline(time.Now().Add(timeout))
		}

		return conn, nil
	}

	// proxy
	if proxyFunc_, ok := options[OPT_PROXY_FUNC]; ok {
		if proxyFunc, ok := proxyFunc_.(func(*http.Request) (int, string, error)); ok {
			transport.Proxy = func(req *http.Request) (*url.URL, error) {
				proxyType, u_, err := proxyFunc(req)
				if err != nil {
					return nil, err
				}

				if proxyType != PROXY_HTTP {
					return nil, fmt.Errorf("only PROXY_HTTP is currently supported")
				}

				u_ = "http://" + u_

				u, err := url.Parse(u_)

				if err != nil {
					return nil, err
				}

				return u, nil
			}
		} else {
			return nil, fmt.Errorf("OPT_PROXY_FUNC is not a desired function")
		}
	} else {
		var proxytype int
		if proxytype_, ok := options[OPT_PROXYTYPE]; ok {
			if proxytype, ok = proxytype_.(int); !ok || proxytype != PROXY_HTTP {
				return nil, fmt.Errorf("OPT_PROXYTYPE must be int, and only PROXY_HTTP is currently supported")
			}
		}

		var proxy string
		if proxy_, ok := options[OPT_PROXY]; ok {
			if proxy, ok = proxy_.(string); !ok {
				return nil, fmt.Errorf("OPT_PROXY must be string")
			}

			if !strings.Contains(proxy, "://") {
				proxy = "http://" + proxy
			}
			proxyUrl, err := url.Parse(proxy)
			if err != nil {
				return nil, err
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	// TLS
	if unsafe_tls_, found := options[OPT_UNSAFE_TLS]; found {
		var unsafe_tls, _ = unsafe_tls_.(bool)
		var tls_config = transport.TLSClientConfig
		if tls_config == nil {
			tls_config = &tls.Config{}
			transport.TLSClientConfig = tls_config
		}
		tls_config.InsecureSkipVerify = unsafe_tls
	}

	return transport, nil
}

// Prepare a redirect policy.
func prepareRedirect(options map[int]interface{}) (func(req *http.Request, via []*http.Request) error, error) {
	var redirectPolicy func(req *http.Request, via []*http.Request) error

	if redirectPolicy_, ok := options[OPT_REDIRECT_POLICY]; ok {
		if redirectPolicy, ok = redirectPolicy_.(func(*http.Request, []*http.Request) error); !ok {
			return nil, fmt.Errorf("OPT_REDIRECT_POLICY is not a desired function")
		}
	} else {
		var followlocation bool
		if followlocation_, ok := options[OPT_FOLLOWLOCATION]; ok {
			if followlocation, ok = followlocation_.(bool); !ok {
				return nil, fmt.Errorf("OPT_FOLLOWLOCATION must be bool")
			}
		}

		var maxredirs int
		if maxredirs_, ok := options[OPT_MAXREDIRS]; ok {
			if maxredirs, ok = maxredirs_.(int); !ok {
				return nil, fmt.Errorf("OPT_MAXREDIRS must be int")
			}
		}

		redirectPolicy = func(req *http.Request, via []*http.Request) error {
			// no follow
			if !followlocation || maxredirs <= 0 {
				return &Error{
					Code:    ErrRedirectPolicy,
					Message: fmt.Sprintf("redirect not allowed"),
				}
			}

			if len(via) >= maxredirs {
				return &Error{
					Code:    ErrRedirectPolicy,
					Message: fmt.Sprintf("stopped after %d redirects", len(via)),
				}
			}

			last := via[len(via)-1]
			// keep necessary headers
			if useragent := last.Header.Get("User-Agent"); useragent != "" {
				req.Header.Set("User-Agent", useragent)
			}

			return nil
		}
	}

	return redirectPolicy, nil
}

// Prepare a cookie jar.
func prepareJar(options map[int]interface{}) (http.CookieJar, error) {
	var jar http.CookieJar
	var err error
	if optCookieJar_, ok := options[OPT_COOKIEJAR]; ok {
		// is bool
		if optCookieJar, ok := optCookieJar_.(bool); ok {
			// default jar
			if optCookieJar {
				jar, err = cookiejar.New(nil)
				if err != nil {
					return nil, err
				}
			}
		} else if optCookieJar, ok := optCookieJar_.(http.CookieJar); ok {
			jar = optCookieJar
		} else {
			return nil, fmt.Errorf("invalid cookiejar")
		}
	}

	return jar, nil
}

// Create an HTTP client.
func NewHttpClient() *HttpClient {
	c := &HttpClient{
		reuseTransport: true,
		reuseJar:       true,
		lock:           new(sync.Mutex),
	}

	return c
}

// Powerful and easy to use HTTP client.
type HttpClient struct {
	// trace_span
	Span zipkin.Span

	// span context key
	SpanContextKey string

	// Default options of client client.
	options map[int]interface{}

	// Default headers of client client.
	Headers map[string]string

	// Options of current request.
	oneTimeOptions map[int]interface{}

	// Headers of current request.
	oneTimeHeaders map[string]string

	// Cookies of current request.
	oneTimeCookies []*http.Cookie

	// Global transport of client client, might be shared between different
	// requests.
	transport http.RoundTripper

	// Global cookie jar of client client, might be shared between different
	// requests.
	jar http.CookieJar

	// Whether current request should reuse the transport or not.
	reuseTransport bool

	// Whether current request should reuse the cookie jar or not.
	reuseJar bool

	// Make requests of one client concurrent safe.
	lock *sync.Mutex

	withLock bool
}

// Set default options and headers.
func (client *HttpClient) Defaults(defaults Map) *HttpClient {
	options, headers := parseMap(defaults)

	// merge options
	if client.options == nil {
		client.options = options
	} else {
		for k, v := range options {
			client.options[k] = v
		}
	}

	// merge headers
	if client.Headers == nil {
		client.Headers = headers
	} else {
		for k, v := range headers {
			client.Headers[k] = v
		}
	}

	return client
}

// Begin marks the begining of a request, it's necessary for concurrent
// requests.
func (client *HttpClient) Begin() *HttpClient {
	client.lock.Lock()
	client.withLock = true

	return client
}

// Reset the client state so that other requests can begin.
func (client *HttpClient) reset() {
	client.oneTimeOptions = nil
	client.oneTimeHeaders = nil
	client.oneTimeCookies = nil
	client.reuseTransport = true
	client.reuseJar = true

	// nil means the Begin has not been called, asume requests are not
	// concurrent.
	if client.withLock {
		client.withLock = false
		client.lock.Unlock()
	}
}

// Temporarily specify an option of the current request.
func (client *HttpClient) WithOption(k int, v interface{}) *HttpClient {
	if client.oneTimeOptions == nil {
		client.oneTimeOptions = make(map[int]interface{})
	}
	client.oneTimeOptions[k] = v

	// Conditions we cann't reuse the transport.
	if hasOption(k, transportOptions) {
		client.reuseTransport = false
	}

	// Conditions we cann't reuse the cookie jar.
	if hasOption(k, jarOptions) {
		client.reuseJar = false
	}

	return client
}

// Temporarily specify multiple options of the current request.
func (client *HttpClient) WithOptions(m Map) *HttpClient {
	options, _ := parseMap(m)
	for k, v := range options {
		client.WithOption(k, v)
	}

	return client
}

// Temporarily specify a header of the current request.
func (client *HttpClient) WithHeader(k string, v string) *HttpClient {
	if client.oneTimeHeaders == nil {
		client.oneTimeHeaders = make(map[string]string)
	}
	client.oneTimeHeaders[k] = v

	return client
}

// Temporarily specify multiple headers of the current request.
func (client *HttpClient) WithHeaders(m map[string]string) *HttpClient {
	for k, v := range m {
		client.WithHeader(k, v)
	}

	return client
}

func (client *HttpClient) WithTrace(ctx interface{}, spanContextKey, name string, trace *zipkin.Tracer) *HttpClient {
	if trace == nil {
		return client
	}

	var span zipkin.Span
	if gutils.IsNil(ctx) {
		span = trace.StartSpan(name)
	} else if ginCtx, ok := ctx.(*gin.Context); ok && ginCtx != nil {
		if ginCtx.Request != (&http.Request{}) {
			span = trace.StartSpan(name, zipkin.Parent(trace.Extract(b3.ExtractHTTP(copyRequest(ginCtx.Request)))))
		} else {
			span = trace.StartSpan(name)
		}
	} else if spanContext, ok := ctx.(model.SpanContext); ok {
		span = trace.StartSpan(name, zipkin.Parent(spanContext))
	} else {
		span = trace.StartSpan(name)
	}

	client.Span = span
	client.SpanContextKey = spanContextKey

	client.WithOption(OPT_BEFORE_REQUEST_FUNC, func(c *http.Client, r *http.Request, spanR zipkin.Span) {
		injector := b3.InjectHTTP(r)
		err := injector(spanR.Context())
		if err != nil {
			log.Printf("gcal before http err: %s", err.Error())
		}
	})

	client.WithOption(OPT_AFTER_REQUEST_FUNC, func(spanR zipkin.Span) {
		if spanR != nil {
			spanR.Finish()
		}
	})
	return client
}

func (client *HttpClient) WithCommonHeader(appName, appSecretKey string) *HttpClient {
	now := time.Now().Unix()
	appSign := gutils.Md5(fmt.Sprintf("%s%d", appSecretKey, now))
	client.WithHeaders(map[string]string{
		"APP-NAME":  appName,
		"APP-SIGN":  appSign,
		"TIMESTAMP": fmt.Sprintf("%d", now),
	})
	return client
}

// Start a request, and get the response.
//
// Usually we just need the Get and Post method.
func (client *HttpClient) Do(method string, url string, headers map[string]string,
	body io.Reader) (*Response, error) {
	options := mergeOptions(defaultOptions, client.options, client.oneTimeOptions)
	headers = mergeHeaders(client.Headers, headers, client.oneTimeHeaders)
	cookies := client.oneTimeCookies

	var transport http.RoundTripper
	var jar http.CookieJar
	var err error

	// transport
	if client.transport == nil || !client.reuseTransport {
		transport, err = prepareTransport(options)
		if err != nil {
			client.reset()
			return nil, err
		}

		if client.reuseTransport {
			client.transport = transport
		}
	} else {
		transport = client.transport
	}

	// jar
	if client.jar == nil || !client.reuseJar {
		jar, err = prepareJar(options)
		if err != nil {
			client.reset()
			return nil, err
		}

		if client.reuseJar {
			client.jar = jar
		}
	} else {
		jar = client.jar
	}

	// release lock
	client.reset()

	redirect, err := prepareRedirect(options)
	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Transport:     transport,
		CheckRedirect: redirect,
		Jar:           jar,
	}

	req, err := prepareRequest(method, url, headers, body, options)
	if err != nil {
		return nil, err
	}
	if debugEnabled, ok := options[OPT_DEBUG]; ok {
		if debugEnabled.(bool) {
			dump, err := httputil.DumpRequestOut(req, true)
			if err == nil {
				fmt.Printf("%s\n", dump)
			}
		}
	}

	if jar != nil {
		jar.SetCookies(req.URL, cookies)
	} else {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	if ctx, ok := options[OPT_CONTEXT]; ok {
		if c, ok := ctx.(context.Context); ok {
			req = req.WithContext(c)
		}
	}

	if beforeReqFunc, ok := options[OPT_BEFORE_REQUEST_FUNC]; ok {
		if f, ok := beforeReqFunc.(func(c *http.Client, r *http.Request, spanR zipkin.Span)); ok {
			f(c, req, client.Span)
		}
	}

	res, err := c.Do(req)

	if afterReqFunc, ok := options[OPT_AFTER_REQUEST_FUNC]; ok {
		if f, ok := afterReqFunc.(func(spanR zipkin.Span)); ok {
			f(client.Span)
		}
	}

	return &Response{res}, err
}

// The HEAD request
func (client *HttpClient) Head(url string) (*Response, error) {
	return client.Do("HEAD", url, nil, nil)
}

// The GET request
func (client *HttpClient) Get(url string, params ...interface{}) (*Response, error) {
	for _, p := range params {
		url = addParams(url, toUrlValues(p))
	}

	return client.Do("GET", url, nil, nil)
}

// The DELETE request
func (client *HttpClient) Delete(url string, params ...interface{}) (*Response, error) {
	for _, p := range params {
		url = addParams(url, toUrlValues(p))
	}

	return client.Do("DELETE", url, nil, nil)
}

// The POST request
//
// With multipart set to true, the request will be encoded as
// "multipart/form-data".
//
// If any of the params key starts with "@", it is considered as a form file
// (similar to CURL but different).
func (client *HttpClient) Post(url string, params interface{}) (*Response,
	error) {
	t := checkParamsType(params)
	if t == 2 {
		return client.Do("POST", url, nil, toReader(params))
	}

	paramsValues := toUrlValues(params)
	// Post with files should be sent as multipart.
	if checkParamFile(paramsValues) {
		return client.PostMultipart(url, params)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body := strings.NewReader(paramsValues.Encode())

	return client.Do("POST", url, headers, body)
}

// Post with the request encoded as "multipart/form-data".
func (client *HttpClient) PostMultipart(url string, params interface{}) (
	*Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	paramsValues := toUrlValues(params)
	// check files
	for k, v := range paramsValues {
		for _, vv := range v {
			// is file
			if k[0] == '@' {
				err := addFormFile(writer, k[1:], vv)
				if err != nil {
					return nil, err
				}
			} else {
				_ = writer.WriteField(k, vv)
			}
		}
	}
	headers := make(map[string]string)

	headers["Content-Type"] = writer.FormDataContentType()
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return client.Do("POST", url, headers, body)
}

func (client *HttpClient) sendJson(method string, url string, data interface{}) (*Response, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	var body []byte
	switch t := data.(type) {
	case []byte:
		body = t
	case string:
		body = []byte(t)
	default:
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	return client.Do(method, url, headers, bytes.NewReader(body))
}

func (client *HttpClient) PostJson(url string, data interface{}) (*Response, error) {
	return client.sendJson("POST", url, data)
}

// The PUT request
func (client *HttpClient) Put(url string, body io.Reader) (*Response, error) {
	return client.Do("PUT", url, nil, body)
}

// Put json data
func (client *HttpClient) PutJson(url string, data interface{}) (*Response, error) {
	return client.sendJson("PUT", url, data)
}
