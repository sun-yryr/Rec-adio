import unittest
from pathlib import Path

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
        self.assertIsNotNone(config["all"].get("keywords"))
        self.assertIsNotNone(config["all"].get("savedir"))
        self.assertIsNotNone(config["all"].get("Radiko_URL"))

    def test_createSaveDir(self):
        """ファイルの作成
        """
        from lib import functions as f
        # 引数がない場合
        path = f.createSaveDirPath()
        current_dir = Path(__file__.replace("test_func.py", ""))
        savefile_path = current_dir.parent / "savefile"
        self.assertEqual(path, str(savefile_path))
        # 引数がある場合
        savefiles = "savefiles"
        path = f.createSaveDirPath(savefiles)
        self.assertEqual(path, savefiles)
        # ファイルの存在確認
        plib = Path(savefiles)
        self.assertTrue(plib.is_dir())
        plib.rmdir()

    def test_recorded_succsess(self):
        from lib import functions as f
        import subprocess
        testfile = "testfile.m4a"
        cmd = "head -c 1000 /dev/urandom > " + testfile
        subprocess.run(cmd, shell=True)
        self.assertFalse(f.is_recording_succeeded(testfile.replace(".m4a", "")))
        Path(testfile).unlink()
        subprocess.run(cmd.replace("1000", "1500"), shell=True)
        self.assertTrue(f.is_recording_succeeded(testfile.replace(".m4a", "")))
        Path(testfile).unlink()
    
    def test_delete_serial(self):
        # なんかいい感じに書きたい
        from lib import functions as f
        inp = "新日曜名作座　雲上雲下  ［終］（９）"
        out = f.delete_serial(inp)
        self.assertEqual(out, "新日曜名作座　雲上雲下")
        inp = "情報学の技術第5回"
        out = f.delete_serial(inp)
        self.assertEqual(out, "情報学の技術")
        

if __name__ == "__main__":
    unittest.main()
