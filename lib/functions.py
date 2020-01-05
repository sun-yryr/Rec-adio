# -*- coding: utf-8 -*-
import os
import json
import requests
import dropbox
import mysql.connector as sql
import hashlib

def load_configurations():
    ROOT = (__file__.replace("/lib/functions.py", ""))
    if (ROOT == __file__):
        ROOT = ROOT.replace("functions.py", "..")
    Path = ROOT + "/conf/config.json"
    if (os.path.isfile(Path) is False):
        print("config file is not found")
        return None
    f = open(Path, "r")
    tmp = json.load(f)
    return tmp

def createSaveDirPath(path = ""):
    ROOT = (__file__.replace("/lib/functions.py", ""))
    if (ROOT == __file__):
        ROOT = ROOT.replace("functions.py", ".")
    if path == "":
        Path = ROOT + "/savefile"
    else:
        Path = path
    if not os.path.isdir(Path):
        os.makedirs(Path)
    return Path

def createSaveDir(Path):
    if (os.path.isdir(Path) is False):
        os.makedirs(Path)

def is_recording_succeeded(Path):
    m4a_path = Path + ".m4a"
    if (os.path.isfile(m4a_path)):
        size = os.path.getsize(m4a_path)
        print("rec size = " + str(size))
        if (size >= 1024):
            return True
        else:
            return False
    else:
        return False

def recording_successful_toline(title):
    headers = {"Authorization": "Bearer %s" % line_token}
    payload = {"message": "\n"+title+" を録音しました!"}
    requests.post("https://notify-api.line.me/api/notify", headers=headers, data=payload)

def recording_failure_toline(title):
    headers = {"Authorization": "Bearer %s" % line_token}
    payload = {"message": "\n"+title+" の録音に失敗しました"}
    requests.post("https://notify-api.line.me/api/notify", headers=headers, data=payload)

class DBXController():
    hadInit = False
    def __init__(self):
        tmpconf = load_configurations()
        if (tmpconf is None) or (tmpconf["all"]["dbx_token"] == ""):
            return
        self.dbx = dropbox.Dropbox(tmpconf["all"]["dbx_token"])
        self.dbx.users_get_current_account()
        res = self.dbx.files_list_folder('')
        db_list = [d.name for d in res.entries]
        if not "radio" in db_list:
            self.dbx.files_create_folder("radio")
        self.hadInit = True

    def upload(self, title, ft, fileData):
        if not self.hadInit:
            return
        dbx_path = "/radio/" + title
        # dropboxにフォルダを作成する
        res = self.dbx.files_list_folder('/radio')
        db_list = [d.name for d in res.entries]
        if not title in db_list:
            self.dbx.files_create_folder(dbx_path)
        dbx_path += "/" +title + "_" + ft[:12]+ ".m4a"
        self.dbx.files_upload(fileData, dbx_path)
    
    def upload_onsen(self, title, count, fileData):
        if not self.hadInit:
            return
        dbx_path = "/radio/" + title
        # dropboxにフォルダを作成する
        res = self.dbx.files_list_folder('/radio')
        db_list = [d.name for d in res.entries]
        if not title in db_list:
            self.dbx.files_create_folder(dbx_path)
        dbx_path += "/" + title + "#" + count + ".mp3"
        self.dbx.files_upload(fileData, dbx_path)

DropBox = DBXController()



class SwiftController():
    hadInit = False
    containerName = "radio"

    def __init__(self):
        tmpconf = load_configurations()
        if (tmpconf is None) or (tmpconf.get("swift") is None):
            return
        self.username = tmpconf["swift"]["username"]
        self.password = tmpconf["swift"]["password"]
        self.tenantid = tmpconf["swift"]["tenantid"]
        self.identityUrl = tmpconf["swift"]["identityUrl"]
        self.objectStrageUrl = tmpconf["swift"]["objectStrageUrl"]
        self.hadInit = True
        # エラーがあったら初期化中止
        if not self.renewal_token():
            print("login error")
            return
        self.create_container(self.containerName)
    
    def renewal_token(self):
        if not self.hadInit:
            return False
        data = {
            "auth": {
                "passwordCredentials": {
                    "username": self.username,
                    "password": self.password
                },
                "tenantId": self.tenantid
            }
        }
        res = requests.post(self.identityUrl + "/tokens",
                            headers={"Content-Type" : "application/json"},
                            data=json.dumps(data))
        resData = json.loads(res.text)
        if "error" in resData.keys():
            return False
        self.token = resData["access"]["token"]["id"]
        return True

    def create_container(self, containerName, isRenewToken = False):
        if not self.hadInit:
            return False
        if isRenewToken:
            self.renewal_token()
        res = requests.put(self.objectStrageUrl + "/" + containerName,
                            headers={
                                "Content-Type" : "application/json",
                                "X-Auth-Token": self.token,
                                "X-Container-Read": ".r:*"
                            })
        if res.status_code in [200, 201, 204]:
            return True
        else:
            return False
    
    def upload_file(self, filePath):
        if not self.hadInit:
            return False
        self.renewal_token()
        # stationとdatetimeでObjectNameを生成する。md5
        hash = hashlib.md5(filePath.encode('utf-8')).hexdigest()
        Path = self.objectStrageUrl + "/" + self.containerName + "/" + hash
        f = open(filePath, "rb")
        res = requests.put(Path,
                            headers={
                                "Content-Type" : "audio/m4a-latm",  # ここで送信するデータ形式を決める
                                "X-Auth-Token": self.token
                            },
                            data=f.read())
        print(res.status_code)
        return Path

Swift = SwiftController()


class DBController:
    hadInit = False

    def __init__(self):
        tmpconf = load_configurations()
        if (tmpconf is None) or (tmpconf.get("mysql") is None):
            return
        self.conn = sql.connect(
            host = tmpconf["mysql"]["hostname"] or 'localhost',
            port = tmpconf["mysql"]["port"] or '3306',
            user = tmpconf["mysql"]["username"],
            password = tmpconf["mysql"]["password"],
            database = tmpconf["mysql"]["database"]
        )
        self.hadInit = True
    
    def insert(self, title, pfm, timestamp, station, uri, info = ""):
        self.conn.ping(reconnect=True)
        cur = self.conn.cursor()
        s = "INSERT INTO Programs (`title`, `pfm`, `rec-timestamp`, `station`, `uri`, `info`) VALUES ( %s, %s, %s, %s, %s, %s)"
        cur.execute(s, (title, pfm, timestamp, station, uri, info))
        self.conn.commit()
        cur.close()

Mysql = DBController()

if __name__ == "__main__":
    test = SwiftController()
    print(test.upload_file("/Users/sun-mm/Downloads/box.m4a"))
    # test = DBController()
    # test.insert()
