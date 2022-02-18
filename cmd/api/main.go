package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/taroooth/rock-paper-scissors/pb"
	"github.com/taroooth/rock-paper-scissors/service"
)

func main() {
	// 起動するポート番号を指定しています。
	port := 50051
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPCサーバーの生成
	server := grpc.NewServer()
	// 自動生成された関数に、サーバと実際に処理を行うメソッドを実装したハンドラを設定します。
	// protoファイルで定義した`RockPaperScissorsService`に対応しています。
	pb.RegisterRockPaperScissorsServiceServer(server, service.NewRockPaperScissorsService())

	// サーバーリフレクションを有効にしています。
	// 有効にすることでシリアライズせずとも後述する`grpc_cli`で動作確認ができるようになります。
	reflection.Register(server)
	// サーバーを起動
	server.Serve(listenPort)
}
