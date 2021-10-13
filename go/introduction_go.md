# 簡單介紹 go
1. goroutine 是 Go 設計並發執行的協程，channel 是在多個 goroutine 中傳遞資料的通道，
   go 在並發上建議使用 CSP 概念 不要以共享內存的方式來通信，要以通過通信共享內存，go 用 channel 達到通信共享內存，
   go 一樣有提供 sync.Mutex 傳統方式處理並發場景。
2. golang.org/x 為 go 標準包的擴充包並比標準有更寬鬆的標準，golang.org/x/sync 提供常見並發組合如 errgroup 。
3. go 中的 panic 會發生程式中止，用 recover 可以取消 panic 恢復程序繼續執行，
   但 recover 第一無法跨 goroutine 處理，第二 go 認為致命的錯誤是無法用 recover 處理，詳細內容請看連結。
   <br/>
   [Panic and crash cases](https://github.com/go101/go101/wiki/Panic-and-crash-cases)
4. 在 go 關鍵字如 append、delete、len 等 當遇到變數是 nil 操作上的行為等同空值，go 關鍵字已處理掉 nil 的判斷邏輯。

* 目錄
    * [參考連結](#參考連結)
    * [常見錯誤](#常見錯誤)
        * [使用 time 發生洩漏場景](#使用-time-發生洩漏場景)
        * [使用 for range slice](#使用-for-range-slice)
        * [循環內使用 defer](#循環內使用-defer)
        * [recover 必須在 defer 內](#recover-必須在-defer-內)
        * [閉包錯誤引用同一個變數](#閉包錯誤引用同一個變數)
        * [空接口不可用 nil 判斷](#空接口不可用-nil-判斷)
        * [不小心覆蓋變量](#不小心覆蓋變量)
        * [go 用 rune 判斷中文長度](#go-用-rune-判斷中文長度)
        * [不導出的 struct 欄位無法被 encode or decode](#不導出的-struct-欄位無法被-encode-or-decode)
        * [receiver 接收為值，無法修改參數原始值](#receiver-接收為值無法修改參數原始值)
        * [switch 中的 fallthrough 語句](#switch-中的-fallthrough-語句)
        * [select 有無 default 差別](#select-有無-default-差別)
        * [select + 多個 channel 優先級](#select--多個-channel-優先級)
        * [HTTP 超時時間](#http-超時時間)
        * [HTTP body 關閉](#http-body-關閉)
        * [log.Fatal 和 log.Panic 不只是 log](#logfatal-和-logpanic-不只是-log)
        * [更新 map 字段的值是 struct 類型](#更新-map-字段的值是-struct-類型)
        * [map for range 順序不固定](#map-for-range-順序不固定)
        * [nil map 操作寫會 panic](#nil-map-操作寫會-panic)
        * [fmt.Println 發生逃逸](#fmtprintln-發生逃逸)
        * [fmt.Println 打印 nil slice or empty slice](#fmtprintln-打印-nil-slice-or-empty-slice)
        * [程序退出不等 goroutine 運行完畢](#程序退出不等-goroutine-運行完畢)
        * [不同 goroutine 之間不滿足順序一致性內存模型](#不同-goroutine-之間不滿足順序一致性內存模型)
        * [go 內建數據結構的操作並不是並發安全的](#go-內建數據結構的操作並不是並發安全的)

## 參考連結
常見錯誤參考以下連結。
<br/>
[Go避坑指南：这些错误你犯过吗？](https://jishuin.proginn.com/p/763bfbd5b767)
<br/>
[Go 语言常见的坑](https://mp.weixin.qq.com/s?__biz=MzAxMTA4Njc0OQ==&mid=2651444773&idx=3&sn=acb089e94cf3549fae9bcf0bde071eb4&chksm=80bb08d7b7cc81c11f61cae49ab4b90c76e43f42d0bceb8c5966d6fbda7134e6fa7b6b421c73&scene=21#wechat_redirect)
<br/>
[Go 并发：一些有趣的现象和要避开的 “坑”](https://segmentfault.com/a/1190000039299143)
<br/>
[Go语言在select语句中实现优先级](https://www.liwenzhou.com/posts/Go/priority_in_go_select/)
<br/>
[附录A：Go语言常见坑](https://chai2010.cn/advanced-go-programming-book/appendix/appendix-a-trap.html)
<br/>
[Golang 新手可能会踩的 50 个坑](https://www.techug.com/post/50-shades-of-go-traps-gotchas-and-common-mistakes__trashed.html)
<br/>
[彻底搞懂-Go逃逸分析](http://www.xiaot123.com/go-e5tbb)
<br/>
[用 Go struct 不能犯的一个低级错误！](https://segmentfault.com/a/1190000040189963?sort=newest)

### 常見錯誤

#### 使用 time 發生洩漏場景
使用 time.NewTimer or time.NewTicker 回傳的 struct 內都有實作一個 `Stop func`，
請記得用 `defer Stop` 關閉資源，但其中有個差異是如果是 Timer 沒有使用 `Stop func` 底層 GC 運作還是會釋放資源，
因 Timer 是一次性任務時間到後就會從底層堆疊中移除，但 Ticker 是定時器任務時間到後會再持續重複一樣的設置時間，
所以使用 time.NewTicker 沒有呼叫 `defer Stop` 則會造成 goroutine 無法回收，而使用 time.NewTimer 沒有呼叫 `defer Stop` 則需要依賴 GC 回收，
使用 time.NewTimer 如果有呼叫 `defer Stop` 則會更早地進行回收，並閱讀理解上可以更易了解回收的時機。

使用 time.After 實作上等同於 time.NewTimer(d).C ，但會特別說 time.After 是因為很常發生在 select + time.After 使用組合，
每次在 select 內的 time.After 都是會產生一個新的 Timer ，在設置時間未到之前是不會被 GC 回收，容易發生洩漏的場景，
可以改使用 time.NewTimer 實作並調整宣告地方避免造成洩漏。

```
for {
   select {
      case: <-time.After(3 * time.Second)
   }
}
```

#### 使用 for range slice
range 是值複製，在使用 range 中去修改值會期待可以改到原始值，
但因為 range 是值複製所以修改 `v.age = 18` 實質上輸出一樣是 `[{25} {32}]`， 請使用 example 2 寫法才可正確改成。

```
// example 1
type person struct {
  age int
}

testSlice := []person{
  {age: 25},
  {age: 32},
}

for _, v := range testSlice {
  v.age = 18
}

fmt.Println(testSlice) // [{25} {32}]
```

```
// example 2
type person struct {
  age int
}

testSlice := []person{
  {age: 25},
  {age: 32},
}

for i := range testSlice {
  testSlice[i].age = 18
}

fmt.Println(testSlice) // [{18} {18}]
```

#### 循環內使用 defer
defer 是在 func 結束時執行，所以如果是在 for 內執行 defer 是會不停的累積值到 func 結束後才真正做關閉資源動作，
如果有需要在 for 內做 defer 執行可以考慮用 example 2 建立一個 匿名 func 提早執行 defer 操作。

```
// example 1
for {
   row, err := db.Query("SELECt ...")
   if err != nil {
      ...
   }
   
   defer row.Close()
}
```

```
// example 2
for {
   func () {
      row, err := db.Query("SELECt ...")
      if err != nil {
         ...
      }
      defer row.Close()
   }()
}
```

#### recover 必須在 defer 內
recover 是可以中止 panic 流程的機制，但必須放在 defer 內才有功能，因為 panic 發生後會接者執行 defer 流程，
如果執行過程中沒有出現 recover 則程序就真的中斷。

#### 閉包錯誤引用同一個變數
請看 example 1 當閉包引用外部參數一律都是指針方式故結果會是 3, 3, 3 ，如果要解決此問題請參考 example 2 or example 3，
使用局部變量或是帶入 func 參數，當有 func 參數時閉包會抓取的是呼叫時的參數。

```
// example 1
for i := 0; i < 3; i++ {
   defer func() {
      fmt.Println(i) // 3, 3, 3
   }()
}
```

```
// example 2
for i := 0; i < 3; i++ {
   i := i
   defer func() {
      fmt.Println(i) // 2, 1, 0
   }()
}
```

```
// example 3
for i := 0; i < 3; i++ {
   defer func(i int) {
      fmt.Println(i) // 2, 1, 0
   }(i)
}
```

#### 空接口不可用 nil 判斷
請看下方範例程式初始化時做判斷邏輯 `p is nil` 當執行 `p = test1` 後結果變為 `p is non-nil`，
原因是 go 的空接口判斷 nil 條件是，沒有值和類型時才是 nil，當執行 `p = test1` p 已獲得 `test struct` 類型訊息資料，
當操作反射時可以取得 `test struct` 相關 Method 或 Field 等相關邏輯，故空接口判斷 nil 條件是，沒有值和類型時才是 nil。

```
type person interface {
   ShowName(name string)
}

type test struct {
}

func (t *test) ShowName(name string) {
   fmt.Println(name)
}

var p person
if p == nil {
  fmt.Println("p is nil")
} else {
  fmt.Println("p is non-nil")
}

var test1 *test
if test1 == nil {
  fmt.Println("test is nil")
} else {
  fmt.Println("test is non-nil")
}

p = test1
if p == nil {
  fmt.Println("p is nil")
} else {
  fmt.Println("p is non-nil")
}

// result
p is nil
test is nil
p is non-nil
```

#### 不小心覆蓋變量
簡短聲明很好用，但在一個區塊內使用 `x := 2` 編譯器不會報錯的會跟預期結果不一致。

```
 x := 1
 fmt.Println(x)        // 1
 {
     fmt.Println(x)    // 1
     x := 2
     fmt.Println(x)    // 2
 }
 fmt.Println(x)        // 1
```

#### go 用 rune 判斷中文長度
go string 儲存是以 utf-8 儲存，但 str 底層是用 byte 數組實現，如要正確取得中文長度需轉換為 `rune` 或使用 `utf8.RuneCountInString`，
byte 常用來處理 ascii、rune 常用來處理 unicode 或 utf-8 。

```
var str = "hello 你好"

fmt.Println("len(str):", len(str)) // 12

fmt.Println("RuneCountInString:", utf8.RuneCountInString(str)) // 8
fmt.Println("rune:", len([]rune(str))) // 8
```

#### 不導出的 struct 欄位無法被 encode or decode
舉例使用 json.Unmarshal 或 json.Marshal 時針對都是 struct 內大寫的欄位，所以 struct 內小寫欄位是會被忽略的，
並且是依照 json tag 邏輯處理，如果大寫但 `json:"-""` 此 - 符號代表忽略的意思。

#### receiver 接收為值，無法修改參數原始值
請看範例一個是 `(p *person)` 一個是 `(p person)` 當各自操作 setHeight 和 setAge 結果只有 height 改成功，
請看 `fmt.Println(fmt.Sprintf("%p", p))` 結果上 `(p person)` 會產生一個新的值所以是不同的位址。

```
type person struct {
	height int
	age    int
}

func (p *person) setHeight(height int) {
	fmt.Println(fmt.Sprintf("%p", p)) // 0xc00009e000
	p.height = height
}

func (p person) setAge(age int) {
	fmt.Println(fmt.Sprintf("%p", &p)) // 0xc00009e040
	p.age = age
}

p := &person{
  height: 170,
  age:    18,
}
fmt.Println(fmt.Sprintf("%p", p)) // 0xc00009e000

p.setHeight(171)
p.setAge(28)

fmt.Println(p) //&{171 18}
```

#### switch 中的 fallthrough 語句
go 跟其它語言的 switch 不一樣，預設會做 break 處理，因為這是常見的預設邏輯，
如果真得有需要執行後續的 case 場景， 使用 fallthrough 可以達到此效果。

#### select 有無 default 差別
select 是給程式處理監聽多個 channel 方式，當所有 channel 都未完成也無 default 時會等待 channel 其中一個完成在往下執行，
當有 default 時如果所有 channel 都未完成，則會執行 default 並往下執行，實際開發上會搭配 `for {}` 所以容易發生一直執行 default ，
而造成 CPU 被一直搶佔情況，如果業務流程上需要則可以像 example 3 一樣，在 default 使用 time 包功能做適當的阻塞。

```
// example 1
t := time.After(2 * time.Second)
t2 := time.After(1 * time.Second)

select {
case <-t:
  fmt.Println("<-t go")
case <-t2:
  fmt.Println("<-t2 go")
}

fmt.Println("finish")
time.Sleep(10 * time.Second)
```

```
// example 2
t := time.After(2 * time.Second)
t2 := time.After(1 * time.Second)

select {
case <-t:
  fmt.Println("<-t go")
case <-t2:
  fmt.Println("<-t2 go")
default:
  fmt.Println("default")
}

fmt.Println("finish")
time.Sleep(10 * time.Second)
```

```
// example 3
for {
   select {
   case: ...
   case: ...
   default:
      // 如使用 time.NewTimer 避免瘋狂搶佔 CPU
   }
}
```

#### select + 多個 channel 優先級
因為 select 當 case 有多個完成是隨機方式執行的，想要有優先級的邏輯就需要使用兩個 select 和搭配使用 break label 方式，
如果像 example 1 當兩個 channel 未完成會陷入死循環執行，所以像 example 2 方式才可達到優先級功能和避免死循環狀況。

```
// example 1
for {
  select {
  case <-stopCh:
      return
  case job1 := <-ch1:
      fmt.Println(job1)
  default:
      select {
      case job2 := <-ch2:
          fmt.Println(job2)
      default:
      }
  }
}
```

```
// example 2
for {
  select {
  case <-stopCh:
      return
  case job1 := <-ch1:
      fmt.Println(job1)
  case job2 := <-ch2:
  priority:
      for {
          select {
          case job1 := <-ch1:
              fmt.Println(job1)
          default:
              break priority
          }
      }
      fmt.Println(job2)
  }
}
```

#### HTTP 超時時間
go http 標準庫 `http.Get、http.Post` 等，預設是沒有設置超時時間，
為確保每個請求會正確關閉需要設置超時時間如 `&http.Client{Timeout: 10 * time.Minute}` 方式。

#### HTTP body 關閉
需要記得使用 `defer` 關閉 `resp.Body.Close()` 但這部分很容易出錯一般會像 example 1 方式撰寫，
正確方式要跟 example 2 ㄧ樣，因如遇到重新定向的錯誤 resp 和 err 都是 non-nil 所以要先判斷 `resp != nil` 在做 `defer` 邏輯，
另一點是 go 把讀取並丟棄數據的任務給 client 處理，go 1.5 版本 client 必須在關閉之前讀取完 body 內資料才可達到 tcp 重用而非關閉。

```
// example 1
client := &http.Client{Timeout: 10 * time.Minute}
resp, err := client.Get("https://www.google.com.tw/?hl=zh_TW")
if err != nil {
  // ...
}

defer resp.Body.Close()
```

```
// example 2
client := &http.Client{Timeout: 10 * time.Minute}
resp, err := client.Get("https://www.google.com.tw/?hl=zh_TW")
if resp != nil {
   defer func(resp *http.Response) {
      if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
         // ...
      }
      
      if err := resp.Body.Close(); err != nil {
         // ...
      }
   }
}
```

#### log.Fatal 和 log.Panic 不只是 log
標準庫的 log.Fatal 和 log.Panic 除了紀錄訊息之外還會調用 os.Exit 或 panic，此兩個動作都會中斷程序，
另外 os.Exit 執行是不會觸發 defer 的，會造成 defer 失效狀況。

#### 更新 map 字段的值是 struct 類型
當 map 的值是 struct 類型時，是無法直接單個更新和取值後更新後是失敗的如 example 1，
正確更新 map 的 struct 值可使用 example 2 或 example 3 方式。

```
type person struct {
   name string
}
```

```
// example 1
m := map[string]person{
  "p1": {"Tom"},
}

//m["p1"].name = "Jerry" // Cannot assign to m["p1"].name

v := m["p1"]
v.name = "Jerry"

fmt.Println(m) // map[p1:{Tom}]
```

```
// example 2
m := map[string]person{
  "p1": {"Tom"},
}

v := m["p1"]
v.name = "Jerry"

m["p1"] = v

fmt.Println(m) // map[p1:{Jerry}]
```

```
// example 3
m := map[string]*person{
  "p1": {"Tom"},
}

m["p1"].name = "Jerry"

fmt.Println(m)       // map[p1:0xc000010240]
fmt.Println(m["p1"]) // &{Jerry}
```

#### map for range 順序不固定
map 是一個哈希表數據結構，此數據結構設計上就不支持順序性的訪問，所以 go 為了避免大家依賴 map 的順序性 key ，
故意在底層增加隨機性的邏輯確保使用者不會誤用。

#### nil map 操作寫會 panic
因為在 go 中的 slice 為 nil 時用 append 一樣可正常使用，容易誤會在其它數據結構也保證了這點，
當在使用 map 賦值時並沒有特別設計一個關鍵字，而 slice 為 nil 可正常賦值是因為 append 關鍵字處理掉 nil 的邏輯。

```
var m map[string]int
m["one"] = 1        // error: panic: assignment to entry in nil map
```

#### fmt.Println 發生逃逸
當使用 `fmt.Println` 參數為 `a ...interface{}` 因為空接口編譯期間無法判斷參數的具體類型，
故會發生逃逸將值宣告到 heap 上，go 逃逸分析是一個很複雜的知識點會隨者版本持續更新，下方連結「彻底搞懂-Go逃逸分析」是提供逃逸分析的場景與基本原理解說，
實際開發上不太需要刻意確認是否發生意逃逸，因為這是底層機制效能優化上的機制，
會說明 `fmt.Println` 的原因是這是常使用的簡單輸出，可能在撰寫一個小程序要驗證一些想法或是比較臨時處理的方式，
可以看下方連結「用 Go struct 不能犯的一个低级错误！」這篇講了一個場景可讓大家理解有時候因為逃逸會造成期望值是不同的。

[彻底搞懂-Go逃逸分析](http://www.xiaot123.com/go-e5tbb)
<br/>
[用 Go struct 不能犯的一个低级错误！](https://segmentfault.com/a/1190000040189963?sort=newest)

#### fmt.Println 打印 nil slice or empty slice
`fmt.Println` 在打印 nil slice or empty slice 都是顯示 `[]`。

```
// nil slice
a := []string{"A", "B", "C", "D", "E"}
a = nil
fmt.Println(a, len(a), cap(a)) // [] 0 0

// 空的slice
a := []string{"A", "B", "C", "D", "E"}
a = a[:0]
fmt.Println(a, len(a), cap(a)) // [] 0 5
```

#### 程序退出不等 goroutine 運行完畢
go ｍain goroutine 當中有開 goroutine 出去，如果 ｍain goroutine 沒有特別做處理等待其它 goroutine 完成，
而直接結束 ｍain goroutine 就真的程序結束不管其它 goroutine 有沒有完成。

#### 不同 goroutine 之間不滿足順序一致性內存模型
go 保證一個 goroutine 內部執行是有順序性的，如果要控制多個 goroutine 有執行順序可使用 channel 控制或其它同步語法。

```
// 無法控制順序性
var msg string
var done bool

func setup() {
    msg = "hello, world"
    done = true
}

func main() {
    go setup()
    for !done {
    }
    println(msg)
}
```

```
// 用 channel 控制順序
var done = make(chan bool)

func setup() {
    msg = "hello, world"
    done <- true
}

func main() {
    go setup()
    <-done
    println(msg)
}
```

#### go 內建數據結構的操作並不是並發安全的
在 go 中使用 slice、map、array 等數據結構並沒有保證同步安全的，只有使用 sync 包的 func 達到同步安全，
但其它部分如 channel、sync.Map 這些是支持同步安全，在使用 goroutine 時候要注意內部使用的數據結構是否同步安全的。
