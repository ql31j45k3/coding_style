package transaction

import "fmt"

func StartTransaction() {
	// cron job 已處理 recover 功能

	// 測試排程有前置條件邏輯，StartTransaction 需要等 StartOrder 完成，才可執行
	fmt.Println("run StartTransaction")
}
