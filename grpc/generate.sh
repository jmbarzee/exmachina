protoc --proto_path=grpc/ --go_out=plugins=grpc:grpc \
  domain.proto \
  dominion.proto \
  service.proto \
  identity.proto