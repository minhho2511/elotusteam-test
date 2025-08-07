# How to run

```
cd hackathon
cp .env.example .env
docker compose up
```

# Services

- Go-kit (https://github.com/go-kit/kit)
- bun (https://bun.uptrace.dev/)
- Redis (backlist token)
- Pg (database)

# APIs

- Register

```
curl --location 'localhost/user/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "minhho2",
    "password": "123456Aa@"
}'

// Response Example
{
    "meta": {
        "correlation_id": "req-8ea5c2b1-b1a1-4c03-b24f-401876aac74c",
        "code": 200,
        "message": "",
        "time": "2025-08-07 03:50:38"
    },
    "data": {}
}
```

- Login

```
curl --location 'localhost/user/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "minhho2",
    "password": "123456Aa@"
}'

// Response Example
{
    "meta": {
        "correlation_id": "req-eaef8361-6a39-4618-ac7f-5be32b870369",
        "code": 200,
        "message": "",
        "time": "2025-08-07 03:52:08"
    },
    "data": {
        "record": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pbmhobzMiLCJpc3MiOiJlbG90dXMtc3lzdGVtIiwic3ViIjoibWluaGhvMyIsImV4cCI6MTc1NDYyNTEyOCwibmJmIjoxNzU0NTM4NzI4LCJpYXQiOjE3NTQ1Mzg3Mjh9.7HBCMs8M2Jz-sfZwiZ3JUpPe5TqngTVoKy1E1wqrWeY",
            "username": "minhho3"
        }
    }
}
```

- File Upload

```
curl --location 'localhost/file/upload' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pbmhobzIiLCJpc3MiOiJlbG90dXMtc3lzdGVtIiwic3ViIjoibWluaGhvMiIsImV4cCI6MTc1NDYyNDUyNywibmJmIjoxNzU0NTM4MTI3LCJpYXQiOjE3NTQ1MzgxMjd9.aoDar5CCunDw_TmDvBnbURuHPgWcJMBSsQ7TpAFVg88' \
--form 'data=@"/Users/minhho/Desktop/hucau.jpeg"'

// Response Example
{
    "meta": {
        "correlation_id": "req-7594a691-607e-448e-a771-ffedc2afd565",
        "code": 200,
        "message": "",
        "time": "2025-08-07 03:42:48"
    },
    "data": {
        "record": {
            "name": "upload_1754538168896707343.jpeg",
            "file_size": 7527,
            "file_type": "image/jpeg",
            "file_path": "/tmp/upload_1754538168896707343.jpeg",
            "info": "{\"ip_address\":\"172.18.0.1\",\"original_name\":\"hucau.jpeg\",\"referer\":\"\",\"user_agent\":\"PostmanRuntime/7.44.1\"}",
            "upload_by": 2,
            "user": null,
            "created_at": "2025-08-07T03:42:48.902263Z",
            "updated_at": "2025-08-07T03:42:48.902263Z",
            "deleted_at": null
        }
    }
}
```

# Struct

- [cfg](cfg)
  - Load config from env
- [cmd](cmd)
  - [main.go](cmd/main.go): main app
- [internal](internal)
  - [initialization](internal/initialization): init router
  - [kit](internal/kit)
    - [endpoints](internal/kit/endpoints): go-kit endpoint (handler mapping)
    - [services](internal/kit/services): go-kit services (biz services)
    - [transports](internal/kit/transports): go-kit transport (now just support http)
  - [middleware](internal/middleware)
    - [jwt.go](internal/middleware/jwt.go): jwt logic
    - [recovery.go](internal/middleware/recovery.go): panic recovery
    - [trace.go](internal/middleware/trace.go): add trade id
  - [models](internal/models): db model
  - [transforms](internal/transforms): mapping request/response
- [migrations](migrations): migration tool
- [pkgs](pkgs): support pkgs
- [utils](utils): support func
- [.env.example](.env.example): env example file, need to copy to .evn
- [docker-compose.yml](docker-compose.yml)

# Note:
- You should commit and push code for each feature, it will prove you know how to use git. Please do not commit all code in one commit => sorry, I missing this field.
- I stick on requirement, no extend feature.