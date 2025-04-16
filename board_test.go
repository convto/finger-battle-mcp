package main

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	if board.CurrentPlayer != PlayerA {
		t.Errorf("初期プレイヤーが間違っています: %v, expected: %v", board.CurrentPlayer, PlayerA)
	}

	if board.PlayerA.Left != 1 || board.PlayerA.Right != 1 {
		t.Errorf("プレイヤーAの初期状態が間違っています: %v, %v", board.PlayerA.Left, board.PlayerA.Right)
	}

	if board.PlayerB.Left != 1 || board.PlayerB.Right != 1 {
		t.Errorf("プレイヤーBの初期状態が間違っています: %v, %v", board.PlayerB.Left, board.PlayerB.Right)
	}
}

func TestApplyFingerCount(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"通常値", 3, 3},
		{"ちょうど5", 5, 0},
		{"5より大きい", 7, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyFingerCount(tt.input)
			if result != tt.expected {
				t.Errorf("ApplyFingerCount(%d) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSelectOpponent(t *testing.T) {
	// 相手の左手を選択するケース
	t.Run("相手の左手を選択", func(t *testing.T) {
		board := NewBoard()
		board.SelectOpponent(1, true) // 左手で相手の左手を攻撃

		// プレイヤーAの左手(1)でプレイヤーBの左手(1)を攻撃すると、プレイヤーBの左手が2になる
		if board.PlayerB.Left != 2 {
			t.Errorf("攻撃後の相手の左手の値が間違っています: %d, expected: 2", board.PlayerB.Left)
		}

		// ターンが切り替わっている
		if board.CurrentPlayer != PlayerB {
			t.Errorf("攻撃後のターンが切り替わっていません: %v, expected: %v", board.CurrentPlayer, PlayerB)
		}
	})

	// 手が消滅するケース
	t.Run("手が消滅するケース", func(t *testing.T) {
		board := NewBoard()
		board.PlayerA.Left = 4
		board.PlayerB.Left = 1

		board.SelectOpponent(1, true) // 左手で相手の左手を攻撃

		// プレイヤーAの左手(4)でプレイヤーBの左手(1)を攻撃すると、プレイヤーBの左手が0(消滅)になる
		if board.PlayerB.Left != 0 {
			t.Errorf("攻撃後の相手の左手が消滅していません: %d, expected: 0", board.PlayerB.Left)
		}
	})

	// 手が5を超えるケース
	t.Run("手が5を超えるケース", func(t *testing.T) {
		board := NewBoard()
		board.PlayerA.Left = 4
		board.PlayerB.Left = 3

		board.SelectOpponent(1, true) // 左手で相手の左手を攻撃

		// プレイヤーAの左手(4)でプレイヤーBの左手(3)を攻撃すると、プレイヤーBの左手が7→2になる
		if board.PlayerB.Left != 2 {
			t.Errorf("攻撃後の相手の左手の値が間違っています: %d, expected: 2", board.PlayerB.Left)
		}
	})
}

func TestRedistributeFingers(t *testing.T) {
	// 通常の再分配
	t.Run("通常の再分配", func(t *testing.T) {
		board := NewBoard()
		board.PlayerA.Left = 2
		board.PlayerA.Right = 3

		success := board.RedistributeFingers(1, 4)

		if !success {
			t.Errorf("再分配に失敗しました")
		}

		if board.PlayerA.Left != 1 || board.PlayerA.Right != 4 {
			t.Errorf("再分配後の値が間違っています: (%d, %d), expected: (1, 4)", board.PlayerA.Left, board.PlayerA.Right)
		}

		// ターンが切り替わっている
		if board.CurrentPlayer != PlayerB {
			t.Errorf("再分配後のターンが切り替わっていません: %v, expected: %v", board.CurrentPlayer, PlayerB)
		}
	})

	// 単純交換は禁止
	t.Run("単純交換は禁止", func(t *testing.T) {
		board := NewBoard()
		board.PlayerA.Left = 2
		board.PlayerA.Right = 3

		success := board.RedistributeFingers(3, 2)

		if success {
			t.Errorf("単純交換が許可されてしまいました")
		}

		// 値が変わっていない
		if board.PlayerA.Left != 2 || board.PlayerA.Right != 3 {
			t.Errorf("単純交換失敗後に値が変わっています: (%d, %d), expected: (2, 3)", board.PlayerA.Left, board.PlayerA.Right)
		}

		// ターンが切り替わっていない
		if board.CurrentPlayer != PlayerA {
			t.Errorf("単純交換失敗後にターンが切り替わっています: %v, expected: %v", board.CurrentPlayer, PlayerA)
		}
	})

	// 再分配後に手が消滅するケース
	t.Run("再分配後に手が消滅するケース", func(t *testing.T) {
		board := NewBoard()
		board.PlayerA.Left = 2
		board.PlayerA.Right = 3

		success := board.RedistributeFingers(5, 0)

		if !success {
			t.Errorf("再分配に失敗しました")
		}

		// 5になった手が消滅している
		if board.PlayerA.Left != 0 || board.PlayerA.Right != 0 {
			t.Errorf("再分配後の値が間違っています: (%d, %d), expected: (0, 0)", board.PlayerA.Left, board.PlayerA.Right)
		}
	})
}

func TestGameOver(t *testing.T) {
	// 先手が勝つケース
	t.Run("先手が勝つケース", func(t *testing.T) {
		board := NewBoard()
		board.PlayerB.Left = 0
		board.PlayerB.Right = 0

		if !board.IsGameOver() {
			t.Errorf("ゲームが終了していないと判定されました")
		}

		if board.GetWinner() != PlayerA {
			t.Errorf("勝者が間違っています: %v, expected: %v", board.GetWinner(), PlayerA)
		}
	})

	// 後手が勝つケース
	t.Run("後手が勝つケース", func(t *testing.T) {
		board := NewBoard()
		board.PlayerA.Left = 0
		board.PlayerA.Right = 0

		if !board.IsGameOver() {
			t.Errorf("ゲームが終了していないと判定されました")
		}

		if board.GetWinner() != PlayerB {
			t.Errorf("勝者が間違っています: %v, expected: %v", board.GetWinner(), PlayerB)
		}
	})

	// ゲームが継続中のケース
	t.Run("ゲームが継続中のケース", func(t *testing.T) {
		board := NewBoard()

		if board.IsGameOver() {
			t.Errorf("ゲームが終了していると誤判定されました")
		}

		if board.GetWinner() != "" {
			t.Errorf("ゲーム継続中なのに勝者がいます: %v", board.GetWinner())
		}
	})
}
