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

if __name__ == "__main__":
    unittest.main()