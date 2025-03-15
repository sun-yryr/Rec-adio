import os
import json

def load_configurations():
    # Construct the path to config.json relative to the project root
    current_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.abspath(os.path.join(current_dir, ".."))
    config_path = os.path.join(project_root, "conf", "config.json")

    try:
        with open(config_path, "r") as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"Error: config file not found at {config_path}")
        return {} # 空の辞書を返す
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in config file at {config_path}")
        return {} # 空の辞書を返す
