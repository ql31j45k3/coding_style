import os
import sys
import paramiko
from dotenv import load_dotenv

# uploadFile2Remote 本地編譯後檔案上傳至遠端伺服器
def uploadFile2Remote(ip, linuxName, linuxPassword, goFilePath, fileName):
    # 建立 ssh 連線
    transport = paramiko.Transport((ip, 22))
    transport.connect(username=linuxName, password=linuxPassword)
    sftp = paramiko.SFTPClient.from_transport(transport)

    localPath = goFilePath + "/" + fileName
    serverPath = "/srv/layout/api/" + fileName

    # sftp 上傳至遠端伺服器
    sftp.put(localPath, serverPath)
    # 修改檔案權限
    sftp.chmod(serverPath, 0o755)

    transport.close()

# goBuild 編譯本地 go 檔案
def goBuild(goFilePath, fileName):
    # 修改目錄位置
    os.chdir(goFilePath)

    # 執行本地編譯 go build make to linux run file
    cmdGoBuild = "GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -race -o " + fileName + " main.go"

    print(os.popen(cmdGoBuild).readlines())

# goBuildDel 刪除編譯後的檔案
def goBuildDel(goFilePath, fileName):
    # 修改目錄位置
    os.chdir(goFilePath)

    cmdDel = "rm " + fileName

    print(os.popen(cmdDel).readlines())

def serviceStop(ip, linuxName, linuxPassword):
    # 建立 ssh 連線
    transport = paramiko.Transport((ip, 22))
    transport.connect(username=linuxName, password=linuxPassword)
    ssh = paramiko.SSHClient()
    ssh._transport = transport

    # 指令停止服務
    stdin, stdout, stderr = ssh.exec_command("systemctl stop layout-api")
    print("[service stop] stdin: ", stdin)
    print("[service stop] stdout: ", stderr)
    print("[service stop] stdout:", stdout.read().decode())

    ssh.close()

def serviceStart(ip, linuxName, linuxPassword):
    # 建立 ssh 連線
    transport = paramiko.Transport((ip, 22))
    transport.connect(username=linuxName, password=linuxPassword)
    ssh = paramiko.SSHClient()
    ssh._transport = transport

    # 指令服務開始
    stdin, stdout, stderr = ssh.exec_command("systemctl start layout-api")
    print("[service start] stdin: ", stdin)
    print("[service start] stdout: ", stderr)
    print("[service start] stdout:", stdout.read().decode())

    stdin, stdout, stderr = ssh.exec_command("systemctl status layout-api")
    print("[check service status] stdin: ", stdin)
    print("[check service status] stdout: ", stderr)
    print("[check service status] stdout:", stdout.read().decode())

    ssh.close()

# mac 本地建置 go build to dev 環境
if __name__ == '__main__':
    load_dotenv()

    ip = os.getenv("DEPLOY_SERVER_IP")
    linuxName = os.getenv("DEPLOY_SERVER_ACCOUNT")
    linuxPassword = os.getenv("DEPLOY_SERVER_PASSWORD")

    # 單一檔案版本
    fileNameList = ["layoutAPI"]

    # 取得當前目錄的上二層
    fullPath = os.getcwd()
    fullPath = fullPath.split("/")[:-1]
    fullPath = "/".join(fullPath)

    for i in range(len(fileNameList)):
        # 組合成目的地目錄
        goFilePath = fullPath + "/cmd/api"

        goBuild(goFilePath, fileNameList[i])

        serviceStop(ip, linuxName, linuxPassword)
        uploadFile2Remote(ip, linuxName, linuxPassword, goFilePath, fileNameList[i])

        goBuildDel(goFilePath, fileNameList[i])
        print("")
        serviceStart(ip, linuxName, linuxPassword)
