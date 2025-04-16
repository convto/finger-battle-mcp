package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Player は先行か後攻かを表す定数
type Player string

const (
	PlayerA Player = "A" // 先行
	PlayerB Player = "B" // 後攻
)

// Hand は手の状態を表す
type Hand struct {
	Left  int
	Right int
}

// Board はゲームの進行状態を表す
type Board struct {
	CurrentPlayer Player // 現在の手番
	PlayerA       Hand   // 先行プレイヤーの手
	PlayerB       Hand   // 後攻プレイヤーの手
}

// NewBoard は初期盤面を作成する
func NewBoard() *Board {
	return &Board{
		CurrentPlayer: PlayerA,
		PlayerA:       Hand{Left: 1, Right: 1},
		PlayerB:       Hand{Left: 1, Right: 1},
	}
}

// String はボードの状態を文字列で返す
func (b *Board) String() string {
	return fmt.Sprintf("プレイヤーA: 左手=%d, 右手=%d\nプレイヤーB: 左手=%d, 右手=%d\n現在の手番: プレイヤー%s",
		b.PlayerA.Left, b.PlayerA.Right, b.PlayerB.Left, b.PlayerB.Right, b.CurrentPlayer)
}

// GetCurrentPlayerHand は現在のプレイヤーの手を返す
func (b *Board) GetCurrentPlayerHand() *Hand {
	if b.CurrentPlayer == PlayerA {
		return &b.PlayerA
	}
	return &b.PlayerB
}

// GetOpponentHand は相手プレイヤーの手を返す
func (b *Board) GetOpponentHand() *Hand {
	if b.CurrentPlayer == PlayerA {
		return &b.PlayerB
	}
	return &b.PlayerA
}

// IsGameOver はゲームが終了したかどうかを判定する
func (b *Board) IsGameOver() bool {
	return (b.PlayerA.Left == 0 && b.PlayerA.Right == 0) ||
		(b.PlayerB.Left == 0 && b.PlayerB.Right == 0)
}

// GetWinner は勝者を返す
func (b *Board) GetWinner() Player {
	if b.PlayerA.Left == 0 && b.PlayerA.Right == 0 {
		return PlayerB
	}
	if b.PlayerB.Left == 0 && b.PlayerB.Right == 0 {
		return PlayerA
	}
	return ""
}

// ApplyFingerCount は指の数を適用し、5以上の場合の処理を行う
func ApplyFingerCount(count int) int {
	if count == 5 {
		return 0 // 消滅
	} else if count > 5 {
		return count - 5 // 5を引く
	}
	return count
}

// SelectOpponent は相手の手を選択する操作
func (b *Board) SelectOpponent(attackingFinger int, isOpponentLeft bool) {
	currentHand := b.GetCurrentPlayerHand()
	opponentHand := b.GetOpponentHand()

	var attackingValue int
	if attackingFinger == 1 { // 左手で攻撃
		attackingValue = currentHand.Left
	} else if attackingFinger == 2 { // 右手で攻撃
		attackingValue = currentHand.Right
	} else {
		fmt.Println("無効な手の選択です")
		return
	}

	// 攻撃する指の数が0の場合は無効
	if attackingValue == 0 {
		fmt.Println("その手は消滅しています")
		return
	}

	// 相手の手に攻撃値を加算
	if isOpponentLeft {
		opponentHand.Left = ApplyFingerCount(opponentHand.Left + attackingValue)
	} else {
		opponentHand.Right = ApplyFingerCount(opponentHand.Right + attackingValue)
	}

	// 手番交代
	b.SwitchTurn()
}

// RedistributeFingers は自分の手の指の数を再分配する
func (b *Board) RedistributeFingers(newLeft, newRight int) bool {
	currentHand := b.GetCurrentPlayerHand()

	// 再分配前の合計
	total := currentHand.Left + currentHand.Right

	// 再分配後の合計が異なる場合は無効
	if total != newLeft+newRight {
		fmt.Println("指の合計数が変わっています")
		return false
	}

	// 自分の手同士の単純交換は禁止
	if currentHand.Left == newRight && currentHand.Right == newLeft {
		fmt.Println("単純な左右交換は禁止されています")
		return false
	}

	// 再分配後の指の数を適用
	currentHand.Left = ApplyFingerCount(newLeft)
	currentHand.Right = ApplyFingerCount(newRight)

	// 手番交代
	b.SwitchTurn()
	return true
}

// SwitchTurn はターンを交代する
func (b *Board) SwitchTurn() {
	if b.CurrentPlayer == PlayerA {
		b.CurrentPlayer = PlayerB
	} else {
		b.CurrentPlayer = PlayerA
	}
}

func main() {
	board := NewBoard()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("===== フィンガーバトル =====")
	fmt.Println("初期盤面:")
	fmt.Println(board)

	for !board.IsGameOver() {
		fmt.Printf("\nプレイヤー%sのターン\n", board.CurrentPlayer)

		// 対象の選択
		fmt.Println("対象を選択してください:")

		// 相手の手の情報を取得
		opponentHand := board.GetOpponentHand()

		if opponentHand.Left > 0 {
			fmt.Printf("1: [攻撃]相手の左手（%d本）\n", opponentHand.Left)
		} else {
			fmt.Println("1: [攻撃]相手の左手（消滅中）")
		}

		if opponentHand.Right > 0 {
			fmt.Printf("2: [攻撃]相手の右手（%d本）\n", opponentHand.Right)
		} else {
			fmt.Println("2: [攻撃]相手の右手（消滅中）")
		}

		fmt.Printf("3: [移動]自分両手指を分配（右手: %d本, 左手: %d本）\n", board.PlayerA.Right, board.PlayerA.Left)

		fmt.Print("> ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1": // 相手の左手を選択
			currentHand := board.GetCurrentPlayerHand()
			fmt.Println("攻撃する手を選択してください:")

			// 消滅している手は選択できないようにする
			attackValidChoices := make(map[string]bool)

			if currentHand.Left > 0 {
				fmt.Printf("1: 左手（%d本）\n", currentHand.Left)
				attackValidChoices["1"] = true
			} else {
				fmt.Println("1: 左手（消滅中）")
			}

			if currentHand.Right > 0 {
				fmt.Printf("2: 右手（%d本）\n", currentHand.Right)
				attackValidChoices["2"] = true
			} else {
				fmt.Println("2: 右手（消滅中）")
			}

			fmt.Print("> ")
			scanner.Scan()
			attackingHandStr := scanner.Text()

			// 消滅している手での攻撃を防ぐ
			if !attackValidChoices[attackingHandStr] {
				fmt.Println("消滅している手では攻撃できません")
				continue
			}

			attackingHand, _ := strconv.Atoi(attackingHandStr)
			board.SelectOpponent(attackingHand, true)

		case "2": // 相手の右手を選択
			currentHand := board.GetCurrentPlayerHand()
			fmt.Println("攻撃する手を選択してください:")

			// 消滅している手は選択できないようにする
			attackValidChoices := make(map[string]bool)

			if currentHand.Left > 0 {
				fmt.Printf("1: 左手（%d本）\n", currentHand.Left)
				attackValidChoices["1"] = true
			} else {
				fmt.Println("1: 左手（消滅中）")
			}

			if currentHand.Right > 0 {
				fmt.Printf("2: 右手（%d本）\n", currentHand.Right)
				attackValidChoices["2"] = true
			} else {
				fmt.Println("2: 右手（消滅中）")
			}

			fmt.Print("> ")
			scanner.Scan()
			attackingHandStr := scanner.Text()

			// 消滅している手での攻撃を防ぐ
			if !attackValidChoices[attackingHandStr] {
				fmt.Println("消滅している手では攻撃できません")
				continue
			}

			attackingHand, _ := strconv.Atoi(attackingHandStr)
			board.SelectOpponent(attackingHand, false)

		case "3": // 自分の手を選択
			currentHand := board.GetCurrentPlayerHand()
			total := currentHand.Left + currentHand.Right

			fmt.Printf("現在の手: 左手（%d本）、右手（%d本）\n", currentHand.Left, currentHand.Right)
			fmt.Printf("現在の合計: %d本\n", total)
			fmt.Printf("新しい左手の値を入力してください (0-%d):\n", total)
			fmt.Print("> ")
			scanner.Scan()
			newLeft, _ := strconv.Atoi(scanner.Text())

			if newLeft < 0 || newLeft > total {
				fmt.Println("無効な値です")
				continue
			}

			newRight := total - newLeft
			fmt.Printf("左手: %d本、右手: %d本 に再分配します\n", newLeft, newRight)
			success := board.RedistributeFingers(newLeft, newRight)
			if !success {
				continue
			}

		default:
			fmt.Println("無効な選択です。1, 2, または 3 を入力してください。")
			continue
		}

		fmt.Println("\n現在の盤面:")
		fmt.Println(board)
	}

	fmt.Printf("\nゲーム終了! プレイヤー%sの勝利です!\n", board.GetWinner())
}
