## index redis
FT.CREATE user-json-idx on json schema $.username as username TEXT SORTABLE $.companyId as companyId TEXT SORTABLE $.status as status TEXT SORTABLE

## sample data create
```
{"topic":"user-request","proc":"new","data":"{\"fullname\":\"Ahmad Zaki\",\"email\":\"ahmad.zaki@gmail.com\",\"hp\":\"+628782345245\",\"username\":\"ahmad\",\"password\":\"12345\",\"companyId\":\"0\",\"roleId\":\"0\",\"officeId\":\"0\",\"departementId\":\"0\"}","from":"","exp":0}
```

```
{"topic":"user-request","proc":"new","data":"{\"fullname\":\"Fahrizal Khoirianto\",\"email\":\"fahrizal.khoirianto@gmail.com\",\"hp\":\"+6287812730567\",\"username\":\"ical\",\"password\":\"ical123\",\"companyId\":\"0\",\"roleId\":\"0\",\"officeId\":\"0\",\"departementId\":\"0\"}","from":"","exp":0}
```

## sample data delete
```
{"topic":"user-request","proc":"delete","data":"{\"username\":\"ical\"}","from":"","exp":0}
```

## sample data update
```
{"topic":"user-request","proc":"edit","data":"{\"username\":\"ical\",\"fullname\":\"Fahrizal\",\"email\":\"info.fahrizal.khoirianto@gmail.com\",\"companyId\":\"4\"}","from":"","exp":0}
```

```
{"topic":"user-request","proc":"edit","data":"{\"fullname\":\"Fahrizal Khoirianto\",\"email\":\"fahrizal.khoirianto@gmail.com\",\"hp\":\"+6287812730567\",\"username\":\"ical\",\"password\":\"ical123\",\"companyId\":\"0\",\"roleId\":\"0\",\"officeId\":\"0\",\"departementId\":\"0\"}","from":"","exp":0}
```

```
{"topic":"user-request","proc":"edit","data":"{\"username\":\"ical\",\"status\":\"1\"}","from":"","exp":0}
```

## sample search
```
{"topic":"user-request","proc":"search","data":"{\"username\":\"ical\",\"fullname\":\"Fahrizal\",\"email\":\"info.fahrizal.khoirianto@gmail.com\",\"companyId\":\"4\"}","from":"","exp":0}
```