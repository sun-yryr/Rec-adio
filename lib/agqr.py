# -*- coding: utf-8 -*-
import lib.functions as f
import time
import subprocess
import os

def rec(data):
    program_data = data[0]
    wait_start_time = data[1]
    program_ft = data[2]
    SAVEROOT = data[3]
    dbx = data[4]
    dir_path = SAVEROOT + "/" + program_data["title"].replace(" ", "_")
    f.createSaveDir(dir_path)
    dbx_path = "/radio/" + program_data["title"]
    res = dbx.files_list_folder('/radio')
    db_list = [d.name for d in res.entries]
    if not program_data["title"] in db_list:
        dbx.files_create_folder(dbx_path)
    dbx_path += "/" +program_data["title"] + "_" + program_ft[:12]+ ".m4a"

    file_path = dir_path + "/" + program_data["title"].replace(" ", "_") + "_" + program_ft[:12]
    cwd  = ('rtmpdump --rtmp "rtmp://fms1.uniqueradio.jp/" ')
    cwd += ('-a ?rtmp://fms-base1.mitene.ad.jp/agqr/ ')
    cwd += ('-f "WIN 16,0,0,257" ')
    cwd += ('-W http://www.uniqueradio.jp/agplayerf/LIVEPlayer-HD0318.swf ')
    cwd += ('-p http://www.uniqueradio.jp/agplayerf/newplayerf2-win.php ')
    cwd += ('-C B:0 ')
    cwd += ('-y aandg1 ')
    cwd += ('--stop %s ' % str(program_data["dur"]))
    cwd += ('--live -o "%s.flv"' % (file_path))
    time.sleep(wait_start_time)
    #rtmpdumpは時間指定の終了ができるので以下を同期処理にする
    subprocess.run(cwd, shell=True)
    #変換をする
    cwd2 = ('ffmpeg -loglevel error -i "%s.flv" -vn -acodec copy "%s.m4a"' % (file_path, file_path))
    subprocess.run(cwd2, shell=True)
    print("agqr finish!")
    if (f.is_recording_succeeded(file_path+ ".m4a")):
        f.recording_successful_toline(program_data["title"])
        fs = open(file_path+".m4a", "rb")
        dbx.files_upload(fs.read(), dbx_path)
        fs.close()
    else:
        f.recording_failure_toline(program_data["title"])
    os.remove(file_path + ".flv")