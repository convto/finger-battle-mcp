# フィンガーバトル (Finger Battle)

フィンガーバトルは、二人のプレイヤーが指の数を増減させて戦うシンプルなゲームです。

## ゲームのルール

1. 二人のプレイヤー（A: 先行、B: 後攻）が対戦します
2. 各プレイヤーは両手（左手と右手）を使います
3. 初期状態では全員の手の指の数は「1」です
4. 各ターンで以下のいずれかの行動を取ります：
   - 相手の手を選択して攻撃：自分の手の指の数を相手の選んだ手に加算します
   - 自分の手の指の数を再分配：合計値を変えずに左右の指の数を変更します

5. 指の数が5になると、その手は消滅（指の数が0になる）します
6. 指の数が5を超えた場合は、5を引いた数になります（例：7 → 2）
7. 相手の両手を消滅させた（両方とも0にした）プレイヤーが勝ちです

## 実行方法

```
go run main.go
```

## 入力方法

1. 対象の選択:
   - 1: 相手の左手
   - 2: 相手の右手
   - 3: 自分の手（移動）

2. 各選択肢での追加入力:
   - 相手の手を選択した場合: 攻撃する自分の手を選びます（0: 左手, 1: 右手）
   - 自分の手を選択した場合: 新しい左手の値を入力します（右手は自動計算）

## 注意点

- 自分の手の単純な左右交換（例: 2,3 → 3,2）は禁止されています
- 一方の手が消滅している場合でも、分配によって復活させることができます
