#!/bin/sh

#入金データ取得 + 未署名トランザクション作成
wallet -r -m 1

#出金データ取得 + 未署名トランザクション作成
wallet -p -m 1

#内部アカウント間の送金 + 未署名トランザクション作成
wallet -t -m 1 -a receipt -t payment

#署名済トランザクションの送信
wallet -s -m 1

#トランザクションのステータス監視
wallet -n -m 1
