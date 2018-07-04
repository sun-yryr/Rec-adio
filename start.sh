#!/bin/bash
GYO=`ps x | grep run.py | wc -l`
URL="https://notify-api.line.me/api/notify"
echo $GYO
if [ $GYO -eq 1 ]; then
    echo "true 1"
    cd /home/pi/hdd1/dev/slackbot/
    MSG=`cat nohup.out`
    curl -X POST -H 'Authorization: Bearer AmCprNqcCYuboLMCqvNRHD3y78lNUnCkKpJkvgJsYO8' -F "message=$MSG"  $URL
    rm nohup.out
    nohup python run.py &
    cd
else
    echo "not 1"
fi