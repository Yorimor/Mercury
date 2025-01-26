Mercury, messenger god and psychopomp.
https://en.wikipedia.org/wiki/Mercury_(mythology)

Requires a folder called `data` which contains two json files

config.json
```json
{
 "auth": "YOUR CLOUDFLARE API TOKEN",
 "ip": ""
}
```

endpoints.json
```json
[
  {
    "name": "DNS RECORD NAME e.g. ftp.domain.com",
    "zone": "DNS ZONE ID",
    "id": "DNS RECORD ID"
  }
]
```

https://developers.cloudflare.com/api/resources/dns/subresources/records/methods/list/