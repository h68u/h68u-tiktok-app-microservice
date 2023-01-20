# api
goctl api go -api ./services/http/app.api -dir ./services/http -style goZero

# rpc
  # user
  goctl rpc protoc ./services/rpc/user/user.proto --go_out=./services/rpc/user/types --go-grpc_out=./services/rpc/user/types --zrpc_out=./services/rpc/user -style goZero

  # video
  goctl rpc protoc ./services/rpc/video/video.proto --go_out=./services/rpc/video/types --go-grpc_out=./services/rpc/video/types --zrpc_out=./services/rpc/video -style goZero

  # contact
  goctl rpc protoc ./services/rpc/contact/contact.proto --go_out=./services/rpc/contact/types --go-grpc_out=./services/rpc/contact/types --zrpc_out=./services/rpc/contact -style goZero
