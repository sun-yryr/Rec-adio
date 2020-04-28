# -*- coding: utf-8 -*-
import os
import json
import requests
import dropbox
import mysql.connector as sql
from bs4 import BeautifulSoup as BS
import hashlib
import subprocess
import time
import re

def load_configurations():
	ROOT = (__file__.replace("/lib/functions.py", ""))
	if (ROOT == __file__):
		ROOT = ROOT.replace("functions.py", "..")
	Path = ROOT + "/conf/config.json"
	if (os.path.isfile(Path) is False):
		print("config file is not found")
		exit(1)
	f = open(Path, "r")
	tmp = json.load(f)
	f.close()
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
		# print("rec size = " + str(size))
		if (size >= 1024):
			return True
		else:
			return False
	else:
		return False

class LineController:
	hadInit = False
	
	def __init__(self):
		tmpconf = load_configurations()
		if (tmpconf.get("all") is None) or (tmpconf["all"].get("line_token") is None):
			return
		self.token = tmpconf["all"]["line_token"]
		self.hadInit = True

	def recording_successful_toline(self, title):
		headers = {"Authorization": "Bearer %s" % self.token}
		payload = {"message": "\n"+title+" を録音しました!"}
		requests.post("https://notify-api.line.me/api/notify", headers=headers, data=payload)

	def recording_failure_toline(self, title):
		headers = {"Authorization": "Bearer %s" % self.token}
		payload = {"message": "\n"+title+" の録音に失敗しました"}
		requests.post("https://notify-api.line.me/api/notify", headers=headers, data=payload)

LINE = LineController()

def did_record_prog(filePath, title, timestamp):
	if (Mysql.hadInit):
		# DBあり
		res = Mysql.check(title, timestamp)
		return (len(res) != 0)
	else:
		# DBなし
		return os.path.exists(filePath)

# delete serial number words
def delete_serial(Path):
	drm_regex = re.compile(r'（.*?）|［.*?］')
	rtn_message = drm_regex.sub("", Path)
	return rtn_message.strip
#

class DBXController():
	hadInit = False
	def __init__(self):
		tmpconf = load_configurations()
		if (tmpconf.get("all") is None) or (tmpconf["all"].get("dbx_token") is None):
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

# rclone

class RcloneController():
	hadInit = False

	def __init__(self):
		tmpconf = load_configurations()
		if (tmpconf.get("rclone") is None):
			return
		self.rcl    = tmpconf["rclone"]["method"]
		self.outdir = tmpconf["rclone"]["outdir"]
		self.rclop  = tmpconf["rclone"]["options"]

		self.hadInit = True

	def upload(self, save_dir, dist_dir):
		if not self.hadInit:
			return
		time.sleep(5)
		cwd = ('rclone %s %s %s %s' % (self.rcl, save_dir, self.outdir+dist_dir+"/" , self.rclop)) 
		p1 = subprocess.run(cwd.split())

Rclone = RcloneController()

# 

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
        self.objectStorageUrl = tmpconf["swift"]["objectStorageUrl"]
        # エラーがあったら初期化中止
        if not self.renewal_token():
            # print("Swift login failed")
            return
        self.hadInit = True
        self.create_container(self.containerName)
    
    def renewal_token(self):
        data = {
            "auth": {
                "passwordCredentials": {
                    "username": self.username,
                    "password": self.password
                },
                "tenantId": self.tenantid
            }
        }
        try:
            res = requests.post(self.identityUrl + "/tokens",
                                headers={"Content-Type" : "application/json"},
                                data=json.dumps(data))
            resData = json.loads(res.text)
            if "error" in resData.keys():
                return False
            self.token = resData["access"]["token"]["id"]
        except:
            return False
        return True

    def create_container(self, containerName, isRenewToken = False):
        if not self.hadInit:
            return False
        if isRenewToken:
            self.renewal_token()
        res = requests.put(self.objectStorageUrl + "/" + containerName,
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
        # create mp3 file
        (root, ext) = os.path.splitext(filePath)
        if (ext == ".m4a"):
            cmd = 'ffmpeg -loglevel error -i "%s" -vn -c:a libmp3lame "%s"' % (filePath, filePath.replace(".m4a", ".mp3"))
            subprocess.run(cmd, shell=True)
        # stationとdatetimeでObjectNameを生成する。md5
        hash = hashlib.md5(filePath.encode('utf-8')).hexdigest()
        Path = self.objectStorageUrl + "/" + self.containerName + "/" + hash
        f = open(filePath.replace(".m4a", ".mp3"), "rb")
        res = requests.put(Path,
                            headers={
                                "Content-Type" : "audio/mpeg",  # ここで送信するデータ形式を決める
                                "X-Auth-Token": self.token
                            },
                            data=f.read())
        print(res.status_code)
        # delete mp3 file
        if (ext == ".m4a"):
            cmd = 'rm "%s"' % (filePath.replace(".m4a", ".mp3"))
            subprocess.run(cmd, shell=True)
        return Path

Swift = SwiftController()


class DBController:
	hadInit = False

	def __init__(self):
		tmpconf = load_configurations()
		if (tmpconf.get("mysql") is None):
			return
		try:
			self.conn = sql.connect(
				host = tmpconf["mysql"]["hostname"] or 'localhost',
				port = tmpconf["mysql"]["port"] or '3306',
				user = tmpconf["mysql"]["username"],
				password = tmpconf["mysql"]["password"],
				database = tmpconf["mysql"]["database"]
			)
			self.hadInit = True
		except:
			# print("Mysql login failed")
			pass
	def insert(self, title, pfm, timestamp, station, uri, info = ""):
		if (not self.hadInit):
			return
		self.conn.ping(reconnect=True)
		cur = self.conn.cursor()
		s = "INSERT INTO Programs (`title`, `pfm`, `rec-timestamp`, `station`, `uri`, `info`) VALUES ( %s, %s, %s, %s, %s, %s)"
		cur.execute(s, (title, pfm, timestamp, station, uri, self.escape_html(info)))
		self.conn.commit()
		cur.close()

	def escape_html(self, html):
		soup = BS(html, "html.parser")
		for script in soup(["script", "style"]):
			script.extract()
		# get text
		text = soup.get_text()
		# break into lines and remove leading and trailing space on each
		lines = (line.strip() for line in text.splitlines())
		# break multi-headlines into a line each
		chunks = (phrase.strip() for line in lines for phrase in line.split("  "))
		# drop blank lines
		text = '\n'.join(chunk for chunk in chunks if chunk)
		return text
	
	def check(self, title, timestamp):
		if (not self.hadInit):
			return
		self.conn.ping(reconnect=True)
		cur = self.conn.cursor()
		s = "SELECT id FROM Programs WHERE `title` = %s AND `rec-timestamp` = %s"
		cur.execute(s, (title, timestamp))
		return cur.fetchall()


Mysql = DBController()

if __name__ == "__main__":
	s = "<img src='http://www.joqr.co.jp/qr_img/detail/20150928195756.jpg' style=\"max-width: 200px;\"> <br /><br /><br />番組メールアドレス：<br /><a href=\"mailto:mar@joqr.net\">mar@joqr.net</a><br />番組Webページ：<br /><a href=\"http://portal.million-arthurs.com/kairi/radio/\">http://portal.million-arthurs.com/kairi/radio/</a><br /><br />パーソナリティは盗賊アーサーを演じる『佐倉綾音』さん、歌姫アーサーを演じる『内田真礼』さん、そして期待の新人『鈴木亜理沙』さん。<br />番組では「乖離性ミリオンアーサー」の最新情報はもちろん、パーソナリティのここだけでしか聞けない話、ゲストをお招きしてのトークなど盛りだくさんでお送りします。<br /><br />初回＆2回目放送は内田真礼さん＆鈴木亜理沙さんのコンビで、その次の2週を佐倉綾音さん＆鈴木亜理沙さんのコンビで2週毎にパーソナリティがローテーションしていく今までにない斬新な番組となります。<br /><br /><br />twitterハッシュタグは「<a href=\"http://twitter.com/search?q=%23millionradio\">#millionradio</a>」<br />twitterアカウントは「<a href=\"http://twitter.com/joqrpr\">@joqrpr</a>」<br />facebookページは「<a href='http://www.facebook.com/1134joqr'>http://www.facebook.com/1134joqr</a>」<br />"
	print(Mysql.escape_html(s))
	# test.insert()
