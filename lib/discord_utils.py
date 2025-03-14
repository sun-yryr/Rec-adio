import os
import discord
import logging
import asyncio
import json

# ロガーの設定
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

async def send_discord_message(channel_id: int, message: str):
    """
    指定されたチャンネルIDにメッセージを送信する非同期関数

    Args:
        channel_id (int): 送信先のチャンネルID
        message (str): 送信するメッセージ内容

    Returns:
        bool: 送信成功 True / 失敗 False
    """
    try:
        config_path = os.path.join(os.path.dirname(__file__), '../conf/config.json')
        with open(config_path, 'r') as f:
            config = json.load(f)
            DISCORD_TOKEN = config['all']['discord_token']

        # Discordクライアントのインスタンスを作成
        client = discord.Client(intents=discord.Intents.default())
        sent = False

        # ログイン処理
        @client.event
        async def on_ready():
            nonlocal sent
            if not sent:
                try:
                    await client.wait_until_ready()  # クライアントが準備完了するまで待つ
                    # チャンネルオブジェクトを取得
                    channel = client.get_channel(channel_id)

                    # メッセージ送信
                    await channel.send(message)
                    logger.info(f"メッセージ送信成功: channel_id={channel_id}, message={message}")
                    sent = True
                except Exception as e:
                    logger.error(f"メッセージ送信失敗: channel_id={channel_id}, message={message}, error={e}")
                finally:
                    await client.close()

        async with client:
            await client.start(DISCORD_TOKEN)
        return sent

    except FileNotFoundError:
        logger.error("config.jsonが見つかりません。")
        return False
    except Exception as e:
        logger.error(f"Discordクライアント起動エラー: {e}")
        return False
