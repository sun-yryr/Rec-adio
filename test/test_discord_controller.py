import asyncio
import unittest
import logging
from unittest.mock import patch
from lib.discord_utils import send_discord_message
from lib.discord_controller import DiscordController
from lib.config_utils import load_configurations

logger = logging.getLogger(__name__)

class TestDiscordController(unittest.IsolatedAsyncioTestCase):
    def setUp(self):
        config = load_configurations()
        if not config or not config.get('all') or not config['all'].get('discord_channel_id') or not config['all'].get('discord_token'):
            self.skipTest("config.json not found or invalid")
        self.controller = DiscordController()
        try:
            self.channel_id = int(config['all']['discord_channel_id'])
            logger.info(f"TestDiscordController initialized: channel_id={self.channel_id}, hadInit={self.controller.hadInit}")
        except (KeyError, ValueError):
            self.skipTest("Invalid config.json: discord_channel_id or discord_token missing or invalid")
        if not self.controller.hadInit:
            self.skipTest("DiscordController not initialized")


    async def test_recording_successful_todiscord(self):
        title = "テスト番組"
        message = f"\n{title} を録音しました!"
        await self.controller.recording_successful_todiscord(title)
        # 送信成功を確認する処理は不要

    async def test_recording_failure_todiscord(self):
        title = "テスト番組"
        message = f"\n{title} の録音に失敗しました"
        await self.controller.recording_failure_todiscord(title)
        # 送信成功を確認する処理は不要

    async def test_recording_failure_todiscord_no_init(self):
        self.controller = DiscordController()
        self.controller.hadInit = False
        title = "テスト番組"
        message = f"\n{title} の録音に失敗しました"
        await self.controller.recording_failure_todiscord(title)
        # 送信成功を確認する処理は不要
