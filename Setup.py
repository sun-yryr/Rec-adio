#!/usr/bin/python3
# -*- coding: utf-8 -*-
import subprocess as shell
import os
import json
import mysql.connector as sql

def main():
    print("configを生成します")
    print("録音ファイルの保存先を指定してください\nデフォルト ./savefile => ", end="")
    tmp = input()
    conf["all"]["savedir"] = tmp if tmp!="" else ""
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
    conf["swift"] = {}
    print("object storage swift")
    print("tenantid => ", end="")
    tmp = input()
    conf["swift"]["tenantid"] = tmp

    print("username => ", end="")
    tmp = input()
    conf["swift"]["username"] = tmp

    print("password => ", end="")
    tmp = input()
    conf["swift"]["password"] = tmp

    print("identityUrl => ", end="")
    tmp = input()
    conf["swift"]["identityUrl"] = tmp

    print("objectStorageUrl => ", end="")
    tmp = input()
    conf["swift"]["objectStorageUrl"] = tmp

def main_mysql():
    conf["mysql"] = {}
    print("mysql")
    print("hostname\nデフォルト localhost => ", end="")
    tmp = input()
    conf["mysql"]["hostname"] = tmp if tmp!="" else "localhost"

    print("port\nデフォルト 3306 => ", end="")
    tmp = input()
    conf["mysql"]["port"] = tmp if tmp!="" else "3306"

    print("username => ", end="")
    tmp = input()
    if tmp == "":
        print("正しい値を入力してください")
        exit(1)
    conf["mysql"]["username"] = tmp

    print("password => ", end="")
    tmp = input()
    if tmp == "":
        print("正しい値を入力してください")
        exit(1)
    conf["mysql"]["password"] = tmp

    print("database name => ", end="")
    tmp = input()
    if tmp == "":
        print("正しい値を入力してください")
        exit(1)
    conf["mysql"]["database"] = tmp

    print("テーブルを自動生成しますか？(新規インストールの方はyを入力してください)\ny/n => ", end="")
    tmp = input()
    if tmp == "y":
        mysql_create_table()

def mysql_create_table():
    try:
        conn = sql.connect(
            host = conf["mysql"]["hostname"],
            port = conf["mysql"]["port"],
            user = conf["mysql"]["username"],
            password = conf["mysql"]["password"],
            database = conf["mysql"]["database"]
        )
        cur = conn.cursor()
        table = "Programs"
        cur.execute("DROP TABLE IF EXISTS `%s`;", table)
        cur.execute(
            """
            CREATE TABLE IF NOT EXISTS `%s` (
            `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
            `title` text NOT NULL,
            `pfm` text,
            `rec-timestamp` datetime NOT NULL,
            `station` varchar(30) DEFAULT '',
            `uri` text NOT NULL,
            `info` text,
            PRIMARY KEY (`id`)
            ) ENGINE=InnoDB AUTO_INCREMENT=395 DEFAULT CHARSET=utf8;
            """, table)
        )
        cur.close()
        conn.close()
    except:
        print("Mysql setup failed")
    

def main_rclone():
    conf["rclone"] = {}
    print("rclone")
    print("method => ", end="")
    tmp = input()
    conf["rclone"]["method"] = tmp

    print("outdir => ", end="")
    tmp = input()
    conf["rclone"]["outdir"] = tmp

    print("options => ", end="")
    tmp = input()
    conf["rclone"]["options"] = tmp

if __name__ == "__main__":
    conf = {
        "all": {
            "savedir": "",
            "Radiko_URL": "http://radiko.jp/v3/program/today/JP13.xml",
            "keywords": []
        },
    }
    main()
    print("Rec-adioを使うにはffmpegとrtmpdumpが必要です。")