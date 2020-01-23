#!/usr/bin/python3
# -*- coding: utf-8 -*-
import subprocess as shell
import os
from Crypto.PublicKey import RSA
from Crypto import Random


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
    # rtmpdump と ffmpeg のインストール確認 返り値を受け取ってパスかどうか確認する
    rtmpPath = shell.check_output(['which', 'rtmpdump']).decode('utf-8').replace('\n', '')
    ffmpegPath = shell.check_output(['which', 'ffmpeg']).decode('utf-8').replace('\n', '')
    if not (os.path.exists(rtmpPath) and os.path.exists(ffmpegPath)):
        print(pycolor.format(pycolor.RED, 'rtmpdump or ffmpeg not found. please download\n'))
        print('ex: ', end='')
        print(pycolor.format(pycolor.ACCENT, 'apt install rtmpdump ffmpeg\n'))
        exit()
    print('ffmpeg, rtmpdump OK.')
    # # react を立てるか確認
    print('Will you use Rec-adio Web Client? y/n : ', end='')
    ans = input()
    if (ans == 'y'):
        # port の指定
        # print('React access port (default = 2332) : ', end='')
        # port = input()
        # if (port == ''): port = '2332'
        # 公開鍵認証の設定
        os.makedirs('./pem', exist_ok=True)
        random_func = Random.new().read
        rsa = RSA.generate(2048, random_func)
        # 秘密鍵作成
        private_pem = rsa.exportKey().decode('utf-8')
        with open('./pem/private.pem', 'w') as f:
            f.write(private_pem)
        # 公開鍵作成
        public_pem = rsa.publickey().exportKey().decode('utf-8')
        with open('./pem/public.pem', 'w') as f:
            f.write(public_pem)
        # ログイン用パスワードの設定
        # print('please passphrase : ', end='')
        # passphrase = input()
        # f = open('./API/.env', 'w')
        # f.write("PORT = '%s'" % port)
        # f.write("PASSWORD = '%s'" % passphrase)
        # f.close()
    elif (ans == 'n'):
        pass
    else:
        print(pycolor.format(pycolor.RED, 'must input y or n.'))
    # config のコピー
    shell.run(['cp', '-n', './conf/example_config.json', './conf/config.json'])
    print('please edit ./conf/config.json')
    # run.py で起動する旨
    print('\nstart for ', end='')
    print(pycolor.format(pycolor.ACCENT, 'python3 run.py'))

if __name__ == "__main__":
    main()