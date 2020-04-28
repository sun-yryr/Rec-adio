import unittest
from pathlib import Path
from lib.onsen import onsen

class TestOnsen(unittest.TestCase):
    SAVEROOT = "./test_save"
    def setUp(self):
        keywords = []
        self.Onsen = onsen(keywords, self.SAVEROOT)
    
    def tearDown(self):
        self.Onsen = None
    
    def test_init(self):
        import datetime as DT
        self.assertEqual(self.Onsen.reload_date, DT.date.today())
        
    def test_change_keywords(self):
        keywords = [
            "麻倉もも", "夏川椎菜", "雨宮天"
        ]
        self.Onsen.change_keywords(keywords)
        self.assertTrue(self.Onsen.keyword.search("麻倉もも"))
        self.assertFalse(self.Onsen.keyword.search("朝倉もも"))

if __name__ == "__main__":
    unittest.main()