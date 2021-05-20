// Powerful and easy to use http client
package gcal

import "sync"

// The default client for convenience
var defaultClient = &HttpClient{
	reuseTransport: true,
	reuseJar:       true,
	lock:           new(sync.Mutex),
}

//var Defaults = defaultClient.Defaults
var Begin = defaultClient.Begin
var Do = defaultClient.Do
var Get = defaultClient.Get
var Delete = defaultClient.Delete
var Head = defaultClient.Head
var Post = defaultClient.Post
var PostJson = defaultClient.PostJson
var PostMultipart = defaultClient.PostMultipart
var Put = defaultClient.Put
var PutJson = defaultClient.PutJson
var WithOption = defaultClient.WithOption
var WithOptions = defaultClient.WithOptions
var WithHeader = defaultClient.WithHeader
var WithHeaders = defaultClient.WithHeaders
var WithTrace = defaultClient.WithTrace
var WithCommonHeader = defaultClient.WithCommonHeader
