# Mercury

### A small tool for updating your cloudflare DNS with your devices public IP

---

There are likely some security implications with using this to open your home network to access from a fixed domain name.

I would **not** recommend using this yourself, and do your own research into the risks.

---

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
    "name": "DNS RECORD NAME e.g. one.domain.com",
    "zone": "DNS ZONE ID",
    "id": "DNS RECORD ID"
  },
  {
    "name": "DNS RECORD NAME e.g. two.domain.com",
    "zone": "DNS ZONE ID",
    "id": "DNS RECORD ID"
  }
]
```

Add an entry for each endpoint you wish to update

Refer to the below cloudflare docs on how to get your dns IDs

https://developers.cloudflare.com/api/resources/dns/subresources/records/methods/list/

---

Mercury, messenger god.