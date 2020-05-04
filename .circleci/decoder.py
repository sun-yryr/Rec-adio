import os
import subprocess
import base64
code=os.environ.get('VPN_SETTING')
with open("tmp.tar.gz", 'wb') as f:
    f.write(base64.b64decode(code))