@echo off
go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2

pushd .
set PATH=%PATH%;./protobuf
set root=..

set pbSrc=%root%\proto\role
set pbDes=%root%\proto\pbRole
del /q %pbDes%\*.pb.go

protoc  --go_out=plugins=grpc:%pbDes% ^
		--proto_path=%pbSrc% ^
		%pbSrc%/*.proto

cd %pbDes%
go mod init pbRole
popd
