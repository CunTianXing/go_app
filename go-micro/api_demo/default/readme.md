## Usage

Run the micro API

```
micro api
```

Run this example

```
go run default.go
```

Make a GET request to /example/call which will call go.micro.api.example Example.Call

```
curl "http://localhost:8080/example/call?name=john"
```

Make a POST request to /example/foo/bar which will call go.micro.api.example Foo.Bar

```
curl -H 'Content-Type: application/json' -d '{}' http://localhost:8080/example/foo/bar


method:"GET"
path:"/example/call"
header:<key:"Accept" value:<key:"Accept" values:"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8" > >
header:<key:"Accept-Encoding" value:<key:"Accept-Encoding" values:"gzip, deflate, br" > > 
header:<key:"Accept-Language" value:<key:"Accept-Language" values:"en-US,en;q=0.8" > > 
header:<key:"Cache-Control" value:<key:"Cache-Control" values: "max-age=0" > > 
header:<key:"Connection" value:<key:"Connection" values:"keep-alive" > > 
header:<key:"Content-Type" value:<key:"Content-Type" values:"application/x-www-form-urlencoded" > > 
header:<key:"Cookie" value:<key:"Cookie" values:"pgv_pvi=1329774592; _ga=GA1.1.1804466211.1497584551; rock_format=json; ROCK_LANG=zh_cn; last-serviceName=sync_grpc"> >
header:<key:"Host" value:<key:"Host" values:"localhost:9099"> > 
header:<key:"Upgrade-Insecure-Requests" value:<key:"Upgrade-Insecure-Requests" values:"1" > > 
header:<key:"User-Agent" value:<key:"User-Agent" values:"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/59.0.3071.109 Chrome/59.0.3071.109 Safari/537.36" > > header:<key:"X-Forwarded-For" value:<key:"X-Forwarded-For" values:"::1" > > 
get:<key:"name" value:<key:"name" values:"john" > >