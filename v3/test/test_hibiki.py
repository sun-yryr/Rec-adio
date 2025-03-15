import unittest
from pathlib import Path
from lib.hibiki import hibiki

class TestHibiki(unittest.TestCase):
    SAVEROOT = "./test_save"
    def setUp(self):
        keywords = []
        self.Hibiki = hibiki(keywords, self.SAVEROOT)
    
    def tearDown(self):
        self.Hibiki = None
    
    def test_init(self):
        self.assertEqual(self.Hibiki.SAVEROOT, self.SAVEROOT)
        
    def test_change_keywords(self):
        keywords = [
            "麻倉もも", "夏川椎菜", "雨宮天"
        ]
        self.Hibiki.change_keywords(keywords)
        self.assertTrue(self.Hibiki.keyword.search("麻倉もも"))
        self.assertFalse(self.Hibiki.keyword.search("朝倉もも"))

if __name__ == "__main__":
    unittest.main()