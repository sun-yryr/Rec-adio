import asyncio
from lib.discord_utils import send_discord_message
import json
import unittest
import os

class TestDiscordUtils(unittest.IsolatedAsyncioTestCase):
    async def asyncSetUp(self):
        config_path = os.path.join(os.path.dirname(__file__), '../conf/config.json')
        with open(config_path, 'r') as f:
            self.config = json.load(f)
            self.channel_id = self.config['all']['discord_channel_id']

    async def test_send_discord_message(self):
        message = "テストメッセージです。"
        success = await send_discord_message(self.channel_id, message)
        self.assertTrue(success)

if __name__ == '__main__':
    unittest.main()
