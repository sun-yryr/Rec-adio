# -*- coding: utf-8 -*-
import requests
import json
from bs4 import BeautifulSoup
import re

c = re.compile(".+>([^<]+)<.+")

def main():
    url = "https://www.agqr.jp/timetable/streaming-sp.php"
    res = requests.get(url)
    res.encoding = 'utf-8'
    soup = BeautifulSoup(res.text, 'html.parser')
    root = soup.find("table", class_="timetb-ag")
    for tr in root.find_all("tr"):
        if (tr.find("div") is not None):
            print("--------------")
            for td in tr.find_all("td"):
                tmp = td.find("div", class_="time")
                #if (tmp.span is not None): tmp.span.decompose()
                time = tmp.get_text().replace("\n", "")
                title = td.find("div", class_="title-p")
                if (title.a is not None):
                    print(time, title.a.string.replace("\n", ""))
                else:
                    print(time, title.string.replace("\n", ""))
        


if __name__ == "__main__":
    main()