
:: 1_user
goctl api go -api ./services/1_user/api/user.api -dir ./services/1_user/api/ -style goZero
goctl rpc protoc ./services/1_user/rpc/user.proto --go_out=./services/1_user/rpc/types --go-grpc_out=./services/1_user/rpc/types --zrpc_out=./services/1_user/rpc -style goZero

:: 2_video
goctl api go -api ./services/2_video/api/video.api -dir ./services/2_video/api/ -style goZero
goctl rpc protoc ./services/2_video/rpc/video.proto --go_out=./services/2_video/rpc/types --go-grpc_out=./services/2_video/rpc/types --zrpc_out=./services/2_video/rpc -style goZero

:: 3_contact
goctl api go -api ./services/3_contact/api/contact.api -dir ./services/3_contact/api/ -style goZero
goctl rpc protoc ./services/3_contact/rpc/contact.proto --go_out=./services/3_contact/rpc/types --go-grpc_out=./services/3_contact/rpc/types --zrpc_out=./services/3_contact/rpc -style goZero