## Trace url redirection(s)

My first project in golang for tracing HTTP redirection

### Build/install

    $ go install rcheck.go http.go utils.go options.go results.go 

### Basic usage

    $ rcheck google.com
    
    [#1] http://google.com - 78.5741ms
     > Status: 301 Moved Permanently
     > Protocol: HTTP/1.1
    
    Redirected to ...
    [#2] http://www.google.com/ - 92.558ms
     > Status: 200 OK
     > Protocol: HTTP/1.1
    2 redirects(s) done in 171.1321ms

### Flag Options

   - `-i` show HTTP headers (default: false)
   - `-b` show the response body content (default: false)
   - `-H` add header to the request (support multiple values)

### More examples

Show redirects headers:

    $ rcheck -i google.com

Simulate a user agent, like a mobile for example:

    $ rcheck -H "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 10_2_1 like Mac OS X) AppleWebKit/602.4.6 (KHTML, like Gecko) Version/10.0 Mobile/14D27 Safari/602.1" google.com



