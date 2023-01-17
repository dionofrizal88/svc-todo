# svc-todo-gRPC

## Running database MYSQL with setting in /server/config/config.json and redis with setting in /server/db/rdb.go:
**Create two schema user and todo using this command':**
```
CREATE DATABASE todo;
```

```
-- todo.todo definition

CREATE TABLE `todo` (
  `row_id` int(11) NOT NULL AUTO_INCREMENT,
  `id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT 0,
  `title` varchar(100) DEFAULT NULL,
  `status` enum('Done','Progress') DEFAULT NULL,
  `done_at` datetime DEFAULT NULL,
  `uuid` varchar(128) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`row_id`),
  KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;
```
```
-- todo.user definition

CREATE TABLE `user` (
  `row_id` int(11) NOT NULL AUTO_INCREMENT,
  `id` bigint(20) unsigned DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `uuid` varchar(128) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`row_id`),
  KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;
```
**Insert user default':**
```
INSERT INTO `user` VALUES (1,0,'None','b38512b6-90fc-11ed-bfc6-0a0027000010');
```
**Create triger to generate column uuid and id:**
**The triger create for two user and todo schema, this example for user:**
```
CREATE DEFINER=`root`@`localhost` TRIGGER `bi_user` BEFORE INSERT ON `user` FOR EACH ROW BEGIN
	if new.id is null then
    SET new.id = UUID_SHORT();
	end if;
	if new.uuid is null then
		set new.uuid = uuid();
	end if;
END
```

## Install protocol buffer for go
**Using this command':**
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
**Update PATH so that the protoc compiler can find the plugins::**
```
$ export PATH="$PATH:$(go env GOPATH)/bin"
```
**Create pb folder**
**Generate gRPC code using MAKEFILE**
```
$ make -f MakeFile gen-proto
```
**Folder proto will contain grpc.pb.go and pb.go**
**Run this command to update go module:**
```
$ go mod tidy
```
**Running golang in folder /server/main.go:**
**Run the following command:**
```
$ go run main.go
```
**Testing grpc request using postman in 0.0.0.0:50051**
**Login using username admin and password secret to get access-token**