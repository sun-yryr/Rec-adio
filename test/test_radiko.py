import unittest
from pathlib import Path
from lib.radiko import radiko

class TestRadiko(unittest.TestCase):
    def setUp(self):
        self.Radiko = radiko()
    
    def tearDown(self):
        self.Radiko = None
    
    def test_init(self):
        import datetime as DT
        self.assertEqual(self.Radiko.reload_date, DT.date.today())
        self.assertIsNotNone(self.Radiko.program_radiko)
        
    def test_change_keywords(self):
        keywords = [
            "麻倉もも", "夏川椎菜", "雨宮天"
        ]
        self.Radiko.change_keywords(keywords)
        self.assertTrue(self.Radiko.keyword.search("麻倉もも"))
        self.assertFalse(self.Radiko.keyword.search("朝倉もも"))

if __name__ == "__main__":
    unittest.main()