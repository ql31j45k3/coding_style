# api
API 文件，使用 API Blueprint 撰寫，選擇 API Blueprint 原因是有幾點，
可以用 Markdown 語法撰寫、可以用文件做 mock server、可以用文件比對服務是否一致，
以上功能可用第三方工具建立，方便上手與使用。


建立 API Blueprint 指令，以下指令生效是有用 npm 安裝 aglio、drakov

e.g. 檔案名稱 api `aglio -i api.apib -o api.html`

建立 mock server 指令

`drakov -f ./api.md -p 3000`