import unittest
import sys
import pathlib
current_dir = pathlib.Path(__file__).resolve().parent
sys.path.append( str(current_dir) + '/../' )

class TestFunction(unittest.TestCase):
    def test_load_lib_functions(self):
        """lib.functions の読み込み
        """
        from lib import functions as f

    def test_load_configurations(self):
        """config.json の読み込み
        """
        from lib import functions as f
        config = f.load_configurations()
        self.assertIsNotNone(config.get("all"))
        self.assertIsNotNone(config.get("swift"))
        self.assertIsNotNone(config.get("mysql"))
        self.assertIsNotNone(config.get("rclone"))

    def test_createSaveDir(self):
        """ファイルの作成
        """
        from lib import functions as f
        path = f.createSaveDirPath()
        current_dir = __file__.replace("test_func.py", "")
        self.assertEqual(path, current_dir + "../savefile")
        savefiles = "savefiles"
        path = f.createSaveDirPath(savefiles)
        self.assertEqual(path, savefiles)
        # ファイルの存在確認
        plib = pathlib.Path(savefiles)
        self.assertTrue(plib.is_dir())
        plib.rmdir()
        self.assertFalse(plib.is_dir())

    def test_recorded_succsess(self):
        from lib import functions as f
        import subprocess
        testfile = "testfile.m4a"
        cmd = "head -c 1000 /dev/urandom > " + testfile
        subprocess.run(cmd, shell=True)
        self.assertFalse(f.is_recording_succeeded(testfile.replace(".m4a", "")))
        pathlib.Path(testfile).unlink()
        subprocess.run(cmd.replace("1000", "1500"), shell=True)
        self.assertTrue(f.is_recording_succeeded(testfile.replace(".m4a", "")))
        pathlib.Path(testfile).unlink()
        

if __name__ == "__main__":
    unittest.main()