# Serviço de usuário utilizando Go e GRPC

## API

### All
```bash
curl --request GET \
  --url http://127.0.0.1:8081/v1/user
```

### Get
```bash
curl --request GET \
  --url http://127.0.0.1:8081/v1/user/1
```

### Create
```bash
curl --request POST \
  --url http://127.0.0.1:8081/v1/user \
  --header 'content-type: application/json' \
  --data '{
	"user": {
		"id": "1",
		"username": "admin",
		"email": "admin@mail.com",
		"state": 2
	}
}'
```

### Update
```bash
curl --request PUT \
  --url http://127.0.0.1:8081/v1/user/1 \
  --header 'content-type: application/json' \
  --data '{
	"user": {
		"id": "1",
		"username": "admin",
		"email": "admin@mail.com",
		"state": 5
	}
}'
```

### Delete
```bash
curl --request DELETE \
  --url http://127.0.0.1:8081/v1/user/1
```