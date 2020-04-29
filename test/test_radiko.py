import unittest
from pathlib import Path
from lib.radiko import radiko
import datetime as DT
from time import sleep

class TestRadiko(unittest.TestCase):
    def setUp(self):
        self.Radiko = radiko()
    
    def tearDown(self):
        self.Radiko = None
    
    def test_init(self):
        self.assertEqual(self.Radiko.reload_date, DT.date.today())
        self.assertIsNotNone(self.Radiko.program_radiko)
        
    def test_change_keywords(self):
        keywords = [
            "麻倉もも", "夏川椎菜", "雨宮天"
        ]
        self.Radiko.change_keywords(keywords)
        self.assertTrue(self.Radiko.keyword.search("麻倉もも"))
        self.assertFalse(self.Radiko.keyword.search("朝倉もも"))
    
    def test_auth(self):
        token = self.Radiko.authorization()
        self.assertIsNotNone(token)

    def test_rec(self):
        # 録音時に使うのは title と dur だけ
        dummy_title = "test_rec_radiko"
        dummy_data = {
            "station": "QRR",
            "title": dummy_title,
            "ft": "20200429044500",
            "dur": 15,
            "pfm": "test",
            "info": "test"
        }
        token = self.Radiko.authorization()
        print("authtoken....=> ", token)
        savePath = Path().resolve() / "save_r"
        savePath.mkdir(parents=True)
        self.Radiko.rec([dummy_data, 0, token, str(savePath)])
        sleep(2)
        # check file
        fileName = dummy_title + "_" + dummy_data["ft"][:12] + ".m4a"
        path = savePath / dummy_title / fileName
        print(str(path))
        self.assertTrue(path.is_file())
        # remove file
        path.unlink()
        path.parent.rmdir()
        savePath.rmdir()

if __name__ == "__main__":
    unittest.main()