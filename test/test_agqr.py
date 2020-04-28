import unittest
from pathlib import Path
from lib.agqr import agqr

class TestAGQR(unittest.TestCase):
    def setUp(self):
        self.Agqr = agqr()
    
    def tearDown(self):
        self.Agqr = None
    
    def test_init(self):
        import datetime as DT
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
        import datetime as DT
        from pathlib import Path
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
        savePath = Path().resolve() / "save"
        print("aaa", str(savePath))
        savePath.mkdir()
        self.Agqr.rec([dummy_data, 0, str(savePath)])
        # check file
        fileName = dummy_title + "_" + dummy_data["ft"] + ".m4a"
        path = savePath / dummy_title / fileName
        self.assertTrue(path.is_file())
        # remove file
        path.unlink()
        path.parent.rmdir()
        savePath.rmdir()


if __name__ == "__main__":
    unittest.main()