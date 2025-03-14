import unittest
from pathlib import Path
from lib.agqr import agqr
import datetime as DT
from time import sleep

class TestAGQR(unittest.TestCase):
    def setUp(self):
        self.Agqr = agqr()
    
    def tearDown(self):
        self.Agqr = None
    
    def test_init(self):
        self.assertEqual(self.Agqr.reload_date, DT.date.today())
        self.assertIs(type(self.Agqr.program_agqr), list)
    
    def test_change_keywords(self):
        keywords = [
            "麻倉もも", "夏川椎菜", "雨宮天"
        ]
        self.Agqr.change_keywords(keywords)
        self.assertTrue(self.Agqr.keyword.search("麻倉もも"))
        self.assertFalse(self.Agqr.keyword.search("朝倉もも"))
    
    def test_rec(self):
        # 録音時に使うのは title と dur だけ
        dummy_title = "test_rec_agqr"
        dummy_data = {
            "title": dummy_title,
            "ft": "202004290000",
            "DT_ft": DT.datetime.strptime("202004290000", "%Y%m%d%H%M"),
            "to": "202004290000",
            "dur": 0.09,
            "pfm": "test"
        }
        savePath = Path().resolve() / "save_a"
        savePath.mkdir(parents=True)
        self.Agqr.rec([dummy_data, 0, str(savePath)])
        sleep(2)
        # check file
        fileName = dummy_title + "_" + dummy_data["ft"][:12] + ".mp3"
        path = savePath / dummy_title / fileName
        self.assertTrue(path.is_file())
        # remove file
        path.unlink()
        path.parent.rmdir()
        savePath.rmdir()


if __name__ == "__main__":
    unittest.main()