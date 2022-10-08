# demo

## layout 設計
相關參考網址
 - [Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_zh-TW.md)
 - [go-clean-arch](https://github.com/bxcodec/go-clean-arch)
 - [在 Golang 中尝试“干净架构”](https://www.jianshu.com/p/5005143e3f4a)
 - [DAY7 — 奔放的 Golang，卻隱藏著有紀律的架構！ — Clean Architecture 實作篇](https://medium.com/%E9%AB%92%E6%A1%B6%E5%AD%90/day7-%E5%A5%94%E6%94%BE%E7%9A%84-golang-%E5%8D%BB%E9%9A%B1%E8%97%8F%E8%91%97%E6%9C%89%E7%B4%80%E5%BE%8B%E7%9A%84%E6%9E%B6%E6%A7%8B-clean-architecture-%E5%AF%A6%E4%BD%9C%E7%AF%87-dd41610dcde7)
 - [DAY13 — Clean Architecture 的力量！無痛從 Restful API 轉換成 gRPC Server](https://medium.com/%E9%AB%92%E6%A1%B6%E5%AD%90/day13-clean-architecture-%E7%9A%84%E5%8A%9B%E9%87%8F-%E7%84%A1%E7%97%9B%E5%BE%9E-restful-api-%E8%BD%89%E6%8F%9B%E6%88%90-grpc-server-e1b64346aa14)
 - [Go 面向包的设计和架构分层](https://github.com/danceyoung/paper-code/blob/master/package-oriented-design/packageorienteddesign.md)
 - [如何写出优雅的 Go 语言代码](https://draveness.me/golang-101/)


    example-1 為按職責拆分，詳細內容可看「相關參考網址或底下簡略內容」。
    example-2 為按層拆分，詳細內容可看「相關參考網址或底下簡略內容」。

    ./layout_2
    ├── cmd： 對外入口點，依照建置服務建立資料夾，內部應只有 main.go 檔案與 main func 做啟動邏輯。
    │   └── demo： 服務啟動流程入口點。
    ├── configs： 設定檔和讀取設定檔邏輯，包含重新讀取 config 邏輯。
    ├── internal： go 導入特殊限制，只允許父級別(layout)與父級別底下子包(api、build 等)導入，可限制其它專案不可導入。
    │   └── domain： 存放所有層，會使用到的對象及方法，提供給 delivery、repository、usecase 呼叫，
    │   │               這樣所有的依賴關係，都是單向的連結 domain，各個實作邏輯不會有 import 關係，
    │   │               如 usecase 不會 import repository 取對象(model) 或是實作程式。
    │   └── example-1： 依照職責拆分，商業邏輯模組，每個資料夾(student、article) 分類功能。
    │       └── student
    │           └── delivery： 該層將充當演示者。決定數據的呈現方式。
    │               │          可以是 REST API、HTML 文件或 GRPC，無論交付類型如何。
    │               │           該層也將處理來自用戶的輸入並將其發送到用例層。
    │               └── http： gin router 註冊 API，包含處理 API request 參數取值與 userCase 回傳值到 response。
    │           └── repository： 處理連線資料庫，實際 SQL 操作邏輯。
    │           └── usecase： 處理核心業務邏輯處理，不能耦合 gin.Context 參數，
    │                           須轉換內部使用 struct，後續可給 GRPC、WebSocket 等 delivery 複用邏輯。
    │       └── article ...
    │   └── example-2： 按層拆分，商業邏輯模組，每個資料夾(usecase、repository) 底下 student_usecase.go、article_usecase.go 分類功能。
    │       └── delivery： 同上。
    │           └── http： 同上。
    │       └── repository： 同上。
    │       └── usecase： 同上。
    │   └── libs： 存放包裝好的組件。
    │       └── logs： log 相關參數設置 (level、path)，包含切割檔案邏輯。
    │       └── middleware： gin middleware 邏輯。
    │       └── mysql： DB 連線初始化、config 設置。
    │       └── response： 封裝 gin 回傳方法，包含回傳格式。
    │   └── utils： 對資料處理的工具，如 string、int、組裝 SQL、反射轉換 struct 等。
    ├── scripts： 執行腳本區。
    │   └── migrate： 數據庫遷移，參考 https://github.com/golang-migrate/migrate。
    │   └── mysql： docker 版 mysql。

## 相關指令與參數  

layout_2/scripts/mysql
有 docker mysql 可使用

```mysql
# 建立 DB SQL
DROP DATABASE IF EXISTS `demo_2`;
CREATE DATABASE IF NOT EXISTS `demo_2` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

```shell
# ide 設定 program arguments:
--configName="config"
--configPath="/Users/michael_kao/go/demo_2/configs"
```
