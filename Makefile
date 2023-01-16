gen-proto:
	@protoc --proto_path=proto proto/*.proto --go_out=E:/Workspace/Golang/src --go-grpc_out=E:/Workspace/Golang/src

cert:
	@cd server/cert; ./gen.sh; cd ..; cd ..