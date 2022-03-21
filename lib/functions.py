# -*- coding: utf-8 -*-
import os
import json
import requests
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
	drm_regex = re.compile(r'（.*?）|［.*?］|第.*?回')
	rtn_message = drm_regex.sub("", Path)
	return rtn_message.rstrip("_ ") 

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
