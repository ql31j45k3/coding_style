# modules
    模組可分為提供者 (dig.Provide)、調用者 (dig.Invoke)，每個模組對外功能都放在模組底下的 dependency 資料夾，
    避免兩個模組變成循環依賴關係，舉例 A 模組使用 B 模組的 dependency，而 B 模組使用 A 模組的 dependency，
    這樣就變成單向依賴的關係，dependency 不會去 import 其他模組，所以都是其它模組單向的 import dependency。

    整個流程都是由 app/api/start 處理，start 負責呼叫各個模組，並使用 dig 管理模組，
    所以不可以有模組 import start，都是 start 單向依賴各個模組，各個模組之間也是單向依賴各個模組的 dependency。
