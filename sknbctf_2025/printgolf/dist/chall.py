#!/usr/bin/env python3
import string
from collections.abc import __builtins__


flag = "sknb{dummy}"
title = """
            _       _              _  __        _           _ _                       
 _ __  _ __(_)_ __ | |_ __ _  ___ | |/ _|   ___| |__   __ _| | | ___ _ __   __ _  ___ 
| '_ \\| '__| | '_ \\| __/ _` |/ _ \\| | |_   / __| '_ \\ / _` | | |/ _ \\ '_ \\ / _` |/ _ \\
| |_) | |  | | | | | || (_| | (_) | |  _| | (__| | | | (_| | | |  __/ | | | (_| |  __/
| .__/|_|  |_|_| |_|\\__\\__, |\\___/|_|_|    \\___|_| |_|\__,_|_|_|\\___|_| |_|\\__, |\\___|
|_|                    |___/                                               |___/         
"""

print(title)
line = input(">>> ")


for c in line:
    if c in string.ascii_letters + string.digits:
        print("Invalid character")
        exit(0)

if len(line) > 8:
    print("Too long")
    exit(0)


bi = __builtins__
del bi["help"]

try:
    eval(line, {"__builtins__": bi}, locals())
except Exception:
    pass
except:
    raise Exception()
