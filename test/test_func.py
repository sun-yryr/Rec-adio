import unittest
from lib.config_utils import load_configurations
import os

class TestLoadConfigurations(unittest.TestCase):
    def test_load_configurations_file_exists(self):
        # config.jsonが存在することを確認
        try:
            config = load_configurations()
            print(f"config: {config}")
            self.assertIsNotNone(config)
            self.assertIn("all", config)
            self.assertIn("discord_token", config["all"])
            self.assertIn("discord_channel_id", config["all"])
        except Exception as e:
            self.fail(f"load_configurations() failed: {e}")


if __name__ == '__main__':
    unittest.main()
