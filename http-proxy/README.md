### Example Request
```
GET /test/1 HTTP/1.1
Host: server.py
Connection: close
User-Agent: netcat
Accept-Language: en


```

### Example Response
```
HTTP/1.0 200 OK
Server: BaseHTTP/0.6 Python/3.7.2
Date: Mon, 26 Oct 2020 07:04:54 GMT
Content-Length: 84

{
    "Host": "server.py",
    "User-Agent": "netcat",
    "Accept-Language": "en"
}
```