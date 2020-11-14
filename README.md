# Doodle

## Create and initialize MySQL database

```
docker run --name doodle-db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql
docker exec -i doodle-db mysql -hlocalhost -u root -proot < schema.sql
```

## Generate TLS certificates

```
mkdir tls && cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go -rsa-bits=2048 -host=localhost
cd -
```
