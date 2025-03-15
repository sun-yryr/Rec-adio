import asyncio
import logging
from lib.config_utils import load_configurations
from lib.discord_utils import send_discord_message

logger = logging.getLogger(__name__)

class DiscordController:
    hadInit = False

    def __init__(self):
        tmpconf = load_configurations()
        if not tmpconf or not tmpconf.get("all") or not tmpconf["all"].get("discord_channel_id") or not tmpconf["all"].get("discord_token"):
            self.hadInit = False
        else:
            self.channel_id = int(tmpconf["all"]["discord_channel_id"])
            self.hadInit = True
            logger.info(f"DiscordController initialized: channel_id={self.channel_id}, hadInit={self.hadInit}")

    async def recording_successful_todiscord(self, title):
        logger.info(f"recording_successful_todiscord called: title={title}, hadInit={self.hadInit}, channel_id={self.channel_id}")
        if not self.hadInit:
            return False
        message = f"\n{title} を録音しました!"
        return await send_discord_message(self.channel_id, message)

    async def recording_failure_todiscord(self, title):
        logger.info(f"recording_failure_todiscord called: title={title}, hadInit={self.hadInit}, channel_id={self.channel_id}")
        if not self.hadInit:
            return False
        message = f"\n{title} の録音に失敗しました"
        return await send_discord_message(self.channel_id, message)
