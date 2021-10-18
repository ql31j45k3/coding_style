# cmd

## shutdown
排程服務與 API 服務 均有做監控 kill or kill -2 指令，
都會通知任務中止訊息，並等待任務正常結束與做資源關閉邏輯。

排程服務：均會等排程執行完。

API 服務：srv.Shutdown 有做時間限制邏輯 (超過時間一樣中止 API 請求)。
