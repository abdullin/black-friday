# run this on commit
import json
import subprocess
from pathlib import Path

print("Commit received. Building...Testing...")

output = subprocess.getoutput("go run *.go")



with Path("/tmp/dat2.json").open('r') as f:
    for l in f:
        l = l.strip()
        if not l:
            continue

        print(l)
        print(json.loads(l))
