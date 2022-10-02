# run this on commit
import subprocess

print("Commit ack")

output = subprocess.getoutput("go run *.go")

print(output)