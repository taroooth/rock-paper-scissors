package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/taroooth/rock-paper-scissors/pb"
	"github.com/taroooth/rock-paper-scissors/pkg"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// DBを使わずに対戦結果の履歴を表示できるように構造体にデータを保持しています。
type RockPaperScissorsService struct {
	numberOfGames int32
	numberOfWins  int32
	matchResults  []*pb.MatchResult
}

// `RockPaperScissorsService`を生成し返却するコンストラクタです。
// `RockPaperScissorsService`は`PlayGame`と`ReportMatchResults`メソッドを実装しています。
func NewRockPaperScissorsService() *RockPaperScissorsService {
	return &RockPaperScissorsService{
		numberOfGames: 0,
		numberOfWins:  0,
		matchResults:  make([]*pb.MatchResult, 0),
	}
}

// 自動生成された`rock-paper-scissors_grpc.pb.go`の
// `RockPaperScissorsServiceServer`インターフェースを実装しています。
func (s *RockPaperScissorsService) PlayGame(ctx context.Context, req *pb.PlayRequest) (*pb.PlayResponse, error) {
	if req.HandShapes == pb.HandShapes_HAND_SHAPES_UNKNOWN {
		return nil, status.Errorf(codes.InvalidArgument, "Choose Rock, Paper, or Scissors.")
	}

	// ランダムに１~３の数値を生成し相手の手とし、`HandShapes`のenumに変換しています。
	opponentHandShapes := pkg.EncodeHandShapes(int32(rand.Intn(3) + 1))

	// ジャンケンの勝敗を決めています。
	var result pb.Result
	if req.HandShapes == opponentHandShapes {
		result = pb.Result_DRAW
	} else if (req.HandShapes.Number()-opponentHandShapes.Number()+3)%3 == 1 {
		result = pb.Result_WIN
	} else {
		result = pb.Result_LOSE
	}


	now := time.Now()
	// 自動生成された型を元に対戦結果を生成
	matchResult := &pb.MatchResult{
		YourHandShapes:     req.HandShapes,
		OpponentHandShapes: opponentHandShapes,
		Result:             result,
		CreateTime: &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}

	// 試合数を１増やし、プレイヤーが勝利した場合は勝利数も１増やします。
	s.numberOfGames = s.numberOfGames + 1
	if result == pb.Result_WIN {
		s.numberOfWins = s.numberOfWins + 1
	}
	s.matchResults = append(s.matchResults, matchResult)

	// 自動生成されたレスポンス用のコードを使ってレスポンスを作り返却しています。
	return &pb.PlayResponse{
		MatchResult: matchResult,
	}, nil
}

// 自動生成された`rock-paper-scissors_grpc.pb.go`の
// `RockPaperScissorsServiceServer`インターフェースを実装しています。
func (s *RockPaperScissorsService) ReportMatchResults(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	// 自動生成されたレスポンス用のコードを使ってレスポンスを作り返却しています。
	return &pb.ReportResponse{
		Report: &pb.Report{
			NumberOfGames: s.numberOfGames,
			NumberOfWins:  s.numberOfWins,
			MatchResults:  s.matchResults,
		},
	}, nil
}
