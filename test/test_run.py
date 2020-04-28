import unittest
from pathlib import Path
import subprocess
from time import sleep
from signal import SIGINT

class TestRun(unittest.TestCase):
    def test_run(self):
        cmd = "pipenv run start"
        proc = subprocess.Popen(cmd.split())
        sleep(15)
        if proc.poll() is not None:
            self.assertIsNone(proc.poll())
        proc.send_signal(SIGINT)
        sleep(2)
        proc.kill()

if __name__ == "__main__":
    unittest.main()