# coding style
* 目錄
    * [參考連結](#參考連結)
    * [檢驗工具](#檢驗工具)
    * [目錄結構](#目錄結構)
    * [指南](#指南)
        * [interface 驗證技巧](#interface-驗證技巧)
        * [零值 Mutex](#零值-mutex)
        * [返回 slice and map](#返回-slice-and-map)
        * [使用 defer 釋放資源](#使用-defer-釋放資源)
        * [goroutine](#goroutine)
        * [使用者決定是否並發](#使用者決定是否並發)
        * [channel size](#channel-size)
        * [channel nil](#channel-nil)
        * [channel struct{}](#channel-struct)
        * [使用 time.Duration](#使用-timeduration)
        * [錯誤處理](#錯誤處理)
        * [錯誤包裝](#錯誤包裝)
        * [使用 panic 時機](#使用-panic-時機)
        * [類型斷言失敗](#類型斷言失敗)
        * [避免使用全域變數](#避免使用全域變數)
        * [避免使用 init](#避免使用-init)
        * [unsafe 包](#unsafe-包)
        * [強制 struct 使用 key: value 初始化](#強制-struct-使用-key-value-初始化)
        * [map 讀取不存在 key](#map-讀取不存在-key)
    * [規範](#規範)
        * [命名風格](#命名風格)
        * [Initialisms (擷取 Go Code Review Comments 中文版)](#initialisms-擷取-go-code-review-comments-中文版)
        * [Receiver Names (擷取 Go Code Review Comments 中文版)](#receiver-names-擷取-go-code-review-comments-中文版)
        * [Variable Names (擷取 Go Code Review Comments 中文版)](#variable-names-擷取-go-code-review-comments-中文版)
        * [interface 命名 (擷取 Effective Go 中文版)](#interface-命名-擷取-effective-go-中文版)
        * [文件命名](#文件命名)
        * [包命名](#包命名)
        * [模組拆分](#模組拆分)
        * [package const and var](#package-const-and-var)
        * [import . 使用時機](#import--使用時機)
        * [context.Context](#contextcontext)
        * [Indent Error Flow](#indent-error-flow)
        * [注意 goroutine 生命週期](#注意-goroutine-生命週期)
    * [性能](#性能)
        * [指定切片容量](#指定切片容量)

## 參考連結
語言指南和規範內制定的規則參考以下連結。
<br/>
[Effective Go 英文版](https://golang.org/doc/effective_go)
<br/>
[Effective Go 中文版](https://go-zh.org/doc/effective_go.html)
<br/>
[Go Code Review Comments 英文版](https://github.com/golang/go/wiki/CodeReviewComments/5a40ba36d388ff1b8b2dd4c1c3fe820b8313152f)
<br/>
[Go Code Review Comments 中文版](https://github.com/panchengtao/articles/issues/8)
<br/>
[uber-go/guide 英文版](https://github.com/uber-go/guide)
<br/>
[uber-go/guide 中文版](https://github.com/xxjwxc/uber_go_guide_cn#%E6%8E%A5%E6%94%B6-slices-%E5%92%8C-maps)
<br/>
[cristaloleg/go-advice 英文版](https://github.com/cristaloleg/go-advice)
<br/>
[cristaloleg/go-advice 中文版](https://github.com/cristaloleg/go-advice/blob/master/README_ZH.md)
<br/>
[Go 开发中的十大常见陷阱[译]](https://tomotoes.com/blog/the-top-10-most-common-mistakes-ive-seen-in-go-projects/)
<br/>
[[go-pkg] context package](https://pjchender.blogspot.com/2020/11/go-pkg-context-package.html)
<br/>
[golang-standards/project-layout](https://github.com/golang-standards/project-layout/blob/master/README_zh-TW.md)
<br/>
[[Golang] 錯誤處理 error handling](https://pjchender.dev/golang/error-handling/)
<br/>
[如何写出优雅的 Go 语言代码](https://draveness.me/golang-101/)

### 檢驗工具
    以下工具應在 IDE 存擋時自動觸發和集成在 CI 流程。
    - gofmt
    - goimports
    - golangci-lint

### 目錄結構
參考 project-layout ，但此 repo 並不是官方標準目前尚未有官方版，此 repo 成為開源社群參考對象。

[golang-standards/project-layout](https://github.com/golang-standards/project-layout/blob/master/README_zh-TW.md)

    ./layout
    ├── api：存放 API 文件或串接第三方 API 文件等。
    ├── build：編譯腳本，如有特殊規則 Drone CI 要求在根目錄可先除外。
    ├── scripts： 相關腳本，如建置、同步 local、dev 環境資料等。
    ├── cmd：對外入口點，依照建置服務建立資料夾，內部應只有 main.go 檔案與 main func 做啟動邏輯。
    ├── configs： 設定檔和讀取設定檔邏輯。
    ├── docs： 在 README.md 相關圖片或存放文件資料。
    ├── internal：go 導入特殊限制，只允許父級別(layout)與父級別底下子包(api、build 等)導入，可限制其它專案不可導入。
    │   └── appapp： 處理初始化流程，因在 cmd/main.go 只能將邏輯寫在 main func 不能拆分檔案，在此處一樣依照建置服務建立資料夾。
    │   └── modules：處理業務邏輯功能模組。
    │   └── utils： 放置共用工具，不包含有商業邏輯的程式。

### 指南

#### interface 驗證技巧
如需在編譯期間檢驗 struct 是否符合 interface 實作，
可在 package var _ 方式，即可編譯期檢驗而且不會對外開放使用。

```
var _ http.Handler = (*Handler)(nil)

type Handler struct {
   // ...
}

func (h *Handler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
)
```

#### 零值 Mutex
無需特定初始化 Mutex 就有 Lock 和 Unlock 效果，初始值等於 go 語言型態的預設值，
另外 Mutex 避免做 copy，因為會造成初始值錯誤，而發生不可預期的死鎖。

```
var mu sync.Mutex
mu.Lock()
```

```
type smap struct {
  sync.Mutex

  data map[string]string
}

func newSMap() *smap {
  return &smap{
    data: make(map[string]string),
  }
}

func (m *smap) Get(k string) string {
  m.Lock()
  defer m.Unlock()

  return m.data[k]
}
```

#### 返回 slice and map
因 slice 和 map 都是一個引用類型，也就是 go 底層 struct 如 slice 結構，
array 是一個引用類型，如果直接回傳原本的 slice 則會發生共用底層 array，
這樣會無法控管那個 func 修改到 slice，原本預期是產生兩個不同 slice 而發生改 A 造成 B 也同步被修改到的耦合性，
除非有在使用 append 關鍵字造成新的擴容才會產生新的 array。

```
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

```
type Stats struct {
  mu sync.Mutex

  counters map[string]int
}

func (s *Stats) Snapshot() map[string]int {
  s.mu.Lock()
  defer s.mu.Unlock()

  result := make(map[string]int, len(s.counters))
  for k, v := range s.counters {
    result[k] = v
  }
  return result
}

// snapshot 现在是一个拷贝
snapshot := stats.Snapshot()
```

#### 使用 defer 釋放資源
使用 defer 做關閉資源動作可增加程式可讀性，因為在取得資源後就在下行程式做關閉行為，
並不會因為程式中間執行邏輯而發生取得資源與關閉資源程式在頭尾，
並且 defer 是在 func 結束後執行，發生 panic 也一樣會做動作，
注意 defer 是後進先出與 defer 會先計算好 func 參數，只有 func 內部邏輯是在最後才執行。

```
p.Lock()
defer p.Unlock()

if p.count < 10 {
  return p.count
}

p.count++
return p.count
```

#### goroutine
go 程式內只要在 `go func` 此 func 就變為並發的方式執行，
在 `go func()` 使用外部變數的話，會造成逃逸代表 sum 生命週期不是當下 func 執行結束就消失，
而是跟者 `go func()` 生命週期和 sum 在 `go func()` 是會自動變成指針形式調用，
建議使用 example 2 寫法可確保變數生命週期與避免因指針調用形式，而發生非預期的耦合行為。

使用 goroutine 可評估是 CPU 資源任務或是 I/O 資源任務，因場景如要解決 I/O 資源任務開更多 goroutine
也可能卡在 I/O 資源等待上排隊，請勿因為跑得慢就開 goroutine 或許在本機測試上會馬上看到執行效能上升，
但服務是持續執行的， 使用 goroutine 會造成非同步執行增加理解程式執行順序的複雜度、
goroutine 可能未正常關閉而洩漏、控管 goroutine 執行的數量「過多是浪費記憶體」等考慮場景。

```
// example 1
func main() {
    c := make(chan struct{})
    sum := 0
   
    go func() {
        sum = 2
        c <- struct{}{}
    }()
   
    <-c
    fmt.Println(sum) // sum = 2
}
```

```
// example 2
func main() {
    c := make(chan struct{})
    sum := 0
   
    go func(sum int) {
        sum = 2
        c <- struct{}{}
    }(sum)
   
    <-c
    fmt.Println(sum) // sum = 0
}
```

#### 使用者決定是否並發
go 官方建議 lib 不應在 func 內使用 goroutine，因使用者無法控制並發行為、並發的數量 (呼叫一次就幾何倍數方式上升)，
但在標準包的處理 API 請求的 ServeHTTP 就寫跟讀都是各一個 goroutine，請求上的寫跟讀為了提高並發量與效能的確使用 goroutine，
並一個 API 請求每次是寫跟讀都是各一個 goroutine 並不會發生倍數的增加曲線，標準包也有提供關閉的 Close func，
並建議註解上提供此邏輯，可讓使用者知道底層重要運行機制。

#### channel size
因 channel 是宣告後就無法做動態調整 size，請依照業務邏輯宣告合適大小，只要大於 1 以上程式就會非同步方式執行，
而不宣告 size 則為無緩沖 channel 代表讀取跟寫是同步執行，如沒有讀到資料則會阻塞程式。

#### channel nil
必免使用 nil channel 因會發生對 nil channel 寫入會永久阻塞、對 nil channel 讀取資料會永久阻塞、
對 nil 操作 close 會發生 panic、 close channel 兩次以上也會發生 panic、向已關閉 channel 發送數據也會 panic。

#### channel struct{}
當 channel 作為一個信號通知的時候，可以使用 struct{} 一個匿名的 struct 內沒有欄位並不會產生內存，
因為一個信號通知只要 channel 有收到資料即可無需判斷資料的值是多少，用 struct{} 可達到通知效果也避免產生內存。

#### 使用 time.Duration
時間單位使用 time.Duration 而不是 int，標準庫 time 已經設定常用單位 Second、Minute 等，
使用 `10*time.Second` 會更清楚表達時間單位。

#### 錯誤處理
go 目前 error handling 是回傳 error，client 端要判斷是否為 nil，的確會增加程式很多判斷式，
go 社群已有討論新的 error handling 版本，下方有連結解釋 err 的使用與一些技巧， 建議都處理回傳的 error 不要使用 _ 忽略。

[[Golang] 錯誤處理 error handling](https://pjchender.dev/golang/error-handling/)

#### 錯誤包裝
官方在 go 1.13 版本加入 Error Wrapping 功能， 請使用 `fmt.Errorf("%w"")` 與 `errors.Is errors.As` 功能做錯誤的判斷邏輯。

#### 使用 panic 時機
使用 panic 會導致服務中斷，在程式初始化流程內，如重要的資源無法建立成功 (DB、config 等)，
就應讓服務啟動失敗，訊息反應服務是不可運行的狀況，而不是持續運行但實際上無法提供功能，
所以初始化流程處理是一個例外，而其他場景應使用 error handling 讓 client 決定如何處理。

#### 類型斷言失敗
是 go 提供從 interface{} (稱為空接口型態) 轉換為原本的型態功能， 要使用 `t, ok := i.(string)` 方式判斷是否轉型成功，
如使用 `t := i.(string)` 轉型失敗會發生 panic 事件，並這是執行期間才可知道的狀況。

```
t, ok := i.(string)
if !ok {
}
```

#### 避免使用全域變數
全域變數是一個方便的宣告但容易濫用，但會造成不易維護、import 依賴，因用全域變數 client 方會主動取得而有 import 依賴問題，
無法管理是誰調用，在寫商業邏輯的程式應使用 di 工具、參數傳遞方式讓依賴反轉，讓程式包關係是單向呼叫。

如會使用「全域變數」特性只有在 lib 包角度 (spf13/viper、sirupsen/logrus)，
lib 包提供的功能特性是有可能給全專案共同使用的狀況，所以在初始化上就提供了一個全域變數並也有提供 `New func` 方式。

#### 避免使用 init
init 是 go 提供初始化的 func，但初始化 func 容易變成初始化變數，但會造成無法控制 init 初始化順序、無法易閱讀初始化流程，
宣告變數是沒有必要性放在初始化 func，應提供 new func 給專案初始化流程內控制使用， 會使用 init 行為如檢查此包必要性邏輯、
或隱式的做註冊邏輯 (database/sql 包)，另外官方 pprof 包就在 init 註冊網址，此方法在社群有很多討論不一定是適合的。

#### unsafe 包
unsafe 包可以提供讓 go 語言安全指針 (內存安全的保證) 變回像是 c 語言版的指針，
可直接對內存做操作處理，但使用此包 go 官方明說不會保證版本的承諾，此包為 go 黑魔法一般是無需使用。

#### 強制 struct 使用 key: value 初始化
go 提供 `Point {1,1}` 簡易方式宣告，但會依賴 struct 變數欄位的順序，為了強制使用 key: value 初始化，
只有增加 `_ struct{}` 就一定要使用 `Point {X：1，Y：1}` 方式宣告不然編譯不會通過，
主要是因為 _ 是等同小寫作用未導出的，所以其他包宣告使用必須要指定 key 方式初始化，
並請 `_ struct{}` 放在 struct 內第一個才不會讓內存提高 (詳細原因請了解 CPU cache)。

```
type Point struct {
  _    struct{} // to prevent unkeyed literals
  X, Y float64
}
```

#### map 讀取不存在 key
go 讀取 map 不存在的 key 會給型態的預設值，不要用預算值做判斷 key 存不存在，
而使用 `v, ok :=` 方式判斷 ok 來做 key 存不存在的邏輯。

```
m1 = make(map[string]string)
v, ok := m1["key"]
if !ok {
}
```

### 規範

#### 命名風格
go 使用 駝峰式 命名包含變量、const、全域 var，等於在 go 程式內部開發上不應出現 駝峰式 以外風格 (舉例如 json tag 等除外)，
命名邏輯思維 請參考 Variable Names， 請注意 go 針對特定命名 如 Id 不應該使用 駝峰式 因為此為縮寫詞應使用 ID (詳細原因了解 Initialisms)。

interface 命名 go 並沒有像 java 一樣使用前綴 I 方式，而是獨立一個風格 請參考 interface 命名，
但 go 只有建議只有一個方法的情況，故在同份專案需討論一個命名風格對於撰寫商業邏輯的 interface 多個方法下的命名風格。

receiver 命名 請參考 Receiver Names，建議取類型的每個間隔的字母 如 agentSettlement 使用 as ，
但如果只有一個字母時並縮寫跟通用重疊 如 v、k，請取一個單字的另外稱呼名字或是可明白的縮寫或直接一樣 (最下策)，
請勿使用 me、this、self 此語言不是面向對象語言會造成理解上的混肴 (混淆：功能跟其它語言是一樣概念嗎？)。

方法內的 receiver 統一一種全部值類型 (as agentSettlement) 或 全部指針類型 (as *agentSettlement)，
實際開發上建議直接統一使用全部指針類型，因大部分場景 struct 內會有混有需要同步修改欄位值或不需要的，統一可增加閱讀性、增加易理解。

#### Initialisms (擷取 Go Code Review Comments 中文版)
名稱中的單詞是首字母或首字母縮略詞（例如 “URL” 或 “NATO” ）需要具有相同的大小寫規則。
例如，“URL” 應顯示為 “URL” 或 “url” （如 “urlPony” 或 “URLPony” ）， 而不是 “Url”。
舉個例子：ServeHTTP 不是 ServeHttp。對於具有多個初始化 “單詞” 的標識符，也應當顯示為 “xmlHTTPRequest” 或 “XMLHTTPRequest”。
當 “ID” 是 “identifier” 的縮寫時，此規則也適用於 “ID” ，因此請寫 “appID” 而不是“appId”。
由協議緩衝區編譯器生成的代碼不受此規則的約束。人工編寫的代碼比機器編寫的代碼要保持更高的標準。

#### Receiver Names (擷取 Go Code Review Comments 中文版)
方法接收者的名稱應該反映其身份；通常，其類型的一個或兩個字母縮寫就足夠了（例如“client”的“c”或“cl”）。
不要使用通用名稱，例如“me”，“this”或“self”，這是面向對象語言的典型標識符，它們更強調方法而不是函數。
名稱不必像方法論證那樣具有描述性，因為它的作用是顯而易見的，不起任何記錄目的。

名稱可以非常短，因為它幾乎出現在每種類型的每個方法的每一行上；familiarity admits brevity。
使用上也要保持一致：如果你在一個方法中叫將接收器命名為“c”，那麼在其他方法中不要把它命名為“cl”。

#### Variable Names (擷取 Go Code Review Comments 中文版)
Go 中的變量名稱應該短而不是長。對於範圍域中的局部變量尤其如此。例如為 line count 定義 c 變量，為 slice index 定義 i 變量。

基本規則：範圍域中，越晚使用的變量，名稱必須越具有描述性。對於方法接收器，一個或兩個字母就足夠了。
諸如循環索引和讀取器（Reader）之類的公共變量可以是單個字母（i，r）。
更多不尋常的事物和全局變量則需要更具描述性的名稱。

#### interface 命名 (擷取 Effective Go 中文版)
按照約定，只包含一個方法的接口應當以該方法的名稱加上-er後綴來命名，如 Reader、Writer、 Formatter、CloseNotifier 等。

諸如此類的命名有很多，遵循它們及其代表的函數名會讓事情變得簡單。 Read、Write、Close、Flush、 String 等都具有典型的簽名和意義。
為避免衝突，請不要用這些名稱為你的方法命名，除非你明確知道它們的簽名和意義相同。
反之，若你的類型實現了的方法， 與一個眾所周知的類型的方法擁有相同含義，那就使用相同命名。
請將字符串轉換方法命名為 String 而非 ToString。

#### 文件命名
一律使用小寫，官方建議文件命名除了測試檔案 (_test.go)，不應出現 _ 命名，
實際開發上還是會有需要做兩個字母的連結，故建議統一使用 _ 命名，不需文件命名出現其它符號。

#### 包命名
包命名為小寫、短命名但不能導致難以意義 如標準包 io、json、os、sort、sync 等，專案內的包名勿跟標準包重名。

包名一般都是單數，但此情形如標準庫的 bytes、errors、strings 為複數，是為了避免跟預定義的類型衝突，
同樣還有go/types是爲了避免和type關鍵字衝突。

實際上知名開源第三方還是會出現使用 - or _ 符號，所以需要符號時請專案統一使用一個符號，
並導入此包時不要使用預設命名需使用重命名方式，如 `reportAPI "test/report-api"`。

#### 模組拆分
在撰寫業務邏輯可以分為按層拆分或按職責拆分，按層就像是 java controller/order.go、service/order.go、dao/order.go ，
按職責拆分就像是 order 資料夾內有 order/controller.go、order/service.go、order/dao.go 方式，
因 go 語言是一個 package 為命名空間並需要自己維護不同包之間的依賴關係，按職責拆分可以符合 go 語言 package 為命名空間的特性，
在處理 order 功能的邏輯 service 跟處理資料庫操作的 dao 都在同一個 package 內可互相引用無需對外，並如果 order 包要改成微服務獨立出去因底層細節沒有對外，
只要將 interface 內的實作改成 http 等連線方式這樣就不會影響到其它包，如要使用按層拆分需要人控制好依賴關係並要確保沒有其它包不小心使用了其它包的 dao，
這是一件很費力的檢驗也無工具可檢測，並如果要拆分微服務並不一定方便修改。

[如何写出优雅的 Go 语言代码](https://draveness.me/golang-101/)

#### package const and var
請依照 group 區分不同類型如 example 2，不要一個像是 example 1 混用。

```
// example 1
const (
   orderA = ""
   orderB = ""
   
   memberA = ""
   memberB = ""
)
```

```
// example 2
const (
   orderA = ""
   orderB = ""
)

const (   
   memberA = ""
   memberB = ""
)
```

#### import . 使用時機
使用 . 可以在程式內不用前綴包命名即可使用，並可避免循環依賴的情況，
但此功能不是在一般程式內使用的引入方法，有循環依賴代表包的關係是需要進行整理的，
此功能是開發給 _test.go 使用的，因在測試場景下可能會為了測試而發生循環依賴是合理的場景。

#### context.Context
go 1.7 引入標準庫內，主要是可以處理多個 goroutine 層次上可一起取消結束 goroutine，
普通場景只需要一個 channel 做關閉信號通知，但如果是我只關閉某個 父 goroutine + 延伸的子 goroutine 場景
單用 channel 是不易管控的。

重要準則
- 不要把 Context 保存在 struct 中，而是直接當作第一個參數傳入 function 或 goroutine 中，通常會命名為 ctx。
- server 在處理傳進來的請求時應該要建立一個 Context，而使用該 server 的方法則應該要接收 Context 作為參數。
- 雖然函式可以允許傳入 nil Context，但千萬不要這麼做，如果你不確定要用哪個 Context，可以使用 context.TODO。
- 只在 request-scoped data 這種要交換處理資料或 API 的範疇下使用 context Values，不要傳入 optional parameters 到函式中。
- 相同的 Context 可以傳入多個不同的 goroutine 中使用，在多個 goroutines 中同時使用 Context 是安全的（safe）。

#### Indent Error Flow
go 建議優先處理錯誤，讓程式先處理錯誤邏輯，可增加可讀性。

```
// 不建議
if err != nil {
    // error handling
} else {
    // normal code
}
```

```
// 建議
if err != nil {
    // error handling
}

// normal code
```

#### 注意 goroutine 生命週期
使用 goroutine 時需清楚知道它如何退出，因持續使用的 goroutine 在 go 底層機制不會去強制關閉 goroutine，
只會盡量去做調度避免發生 某個 goroutine 取得不到資源餓死。
goroutine 與 channel 是搭配使用的並發模型， 但 goroutine 對 channel 做讀寫操作可能會發生堵塞，
並持續堵塞會發生記憶體持續升高或無法釋放的情況， 此情況 GC 也無能為力此現象稱為 goroutine 洩漏。

### 性能

#### 指定切片容量
使用 `make([]T, length, capacity)` 指定 capacity 做切片的初始化，
在一般業務邏輯的場景可以依照分頁大小或清楚知道 capacity 的大小場景，使用 make 需指定 capacity 做初始化，
如無法判斷 capacity 大小場景，就使用 `make([]T, length)` 讓底層機制依照擴容條件進行擴容。
