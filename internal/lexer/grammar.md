## Context-Free Grammar for HTTP Request

```
POST /api/data HTTP/1.1
Host: localhost:5000
User-Agent: curl/7.88.1
Accept: */*
Content-Type: application/json
X-Request-ID: 12345-random
Content-Length: 47

{"name": "testuser", "age": 30, "active": true}
```

#### Tokens
- `VERB`: POST | GET | PATCH | PUT | DELETE
- `ROUTE`: /, /api, /*
- `PROTOCOL`: HTTP, HTTPS
- `PROTO_VERSION`: 1.1, 2
- `HOST`: localhost:3000, 168.1.2.8:5000
- `ACCEPT`: application/json, */*
- `HEADER`: any key-value pair not explicitly accounted for
- `CONTENT_LENGTH`: N value
- `CONTENT`: any data in any *supported* format

