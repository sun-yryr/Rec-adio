#!/usr/bin/python3
# -*- coding: utf-8 -*-
import subprocess as shell
import os
import json


class pycolor:
    BLACK = '\033[30m'
    RED = '\033[31m'
    GREEN = '\033[32m'
    YELLOW = '\033[33m'
    BLUE = '\033[34m'
    PURPLE = '\033[35m'
    CYAN = '\033[36m'
    WHITE = '\033[37m'
    RETURN = '\033[07m' #反転
    ACCENT = '\033[01m' #強調
    FLASH = '\033[05m' #点滅
    RED_FLASH = '\033[05;41m' #赤背景+点滅
    END = '\033[0m'
    def format(Color, text):
        return Color + text + pycolor.END

def main():
    print("configを生成します")
    print("録音ファイルの保存先を指定してください\nデフォルト ./savefile => ", end="")
    tmp = input()
    if tmp != "":
        conf["all"]["savedir"] = tmp
    print("lineで通知する場合はトークンを入力してください => ", end="")
    tmp = input()
    if tmp != "":
        conf["all"]["line_token"] = tmp
    print("オブジェクトストレージを使いますか？ y/n => ", end="")
    tmp = input()
    if tmp == "y":
        main_objectstorage()
    print("mysqlを使いますか？ y/n => ", end="")
    tmp = input()
    if tmp == "y":
        main_mysql()
    print("rcloneを使いますか？ y/n => ", end="")
    tmp = input()
    if tmp == "y":
        main_rclone()
    print("./conf/config.jsonに保存しました")
    with open("./conf/config.json", "w") as f:
        json.dump(conf, f, ensure_ascii=False)


def main_objectstorage():
    pass

def main_mysql():
    pass

def main_rclone():
    pass

if __name__ == "__main__":
    conf = {
        "all": {
            "savedir": "",
            "Radiko_URL": "http://radiko.jp/v3/program/today/JP13.xml",
            "keywords": []
        },
    }
    main()