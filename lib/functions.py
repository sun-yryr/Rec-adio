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
	return os.path.exists(filePath)

# delete serial number words
def delete_serial(Path):
	drm_regex = re.compile(r'（.*?）|［.*?］|第.*?回')
	rtn_message = drm_regex.sub("", Path)
	return rtn_message.rstrip("_ ") 
