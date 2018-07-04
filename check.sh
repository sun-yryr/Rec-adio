#!/bin/bash

URL="https://notify-api.line.me/api/notify"
MSG="超A&G URL変更"
SIZE=`wc -c rec_radiko.sh |cut -d ' ' -f 5`
if [ $SIZE -lt 500 ]; then
    curl -X POST -H 'Authorization: Bearer AmCprNqcCYuboLMCqvNRHD3y78lNUnCkKpJkvgJsYO8' -F "message=$MSG"  $URL
fi