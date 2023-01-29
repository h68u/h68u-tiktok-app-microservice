cd ../
:: api
goctl api go -api ./service/http/app.api -dir ./service/http -style goZero

:: rpc
  :: user
  goctl rpc protoc ./service/rpc/user/user.proto --go_out=./service/rpc/user/types --go-grpc_out=./service/rpc/user/types --zrpc_out=./service/rpc/user -style goZero

  :: video
  goctl rpc protoc ./service/rpc/video/video.proto --go_out=./service/rpc/video/types --go-grpc_out=./service/rpc/video/types --zrpc_out=./service/rpc/video -style goZero

  :: contact
  goctl rpc protoc ./service/rpc/contact/contact.proto --go_out=./service/rpc/contact/types --go-grpc_out=./service/rpc/contact/types --zrpc_out=./service/rpc/contact -style goZero
