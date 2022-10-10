#!/usr/bin/python3
import requests, sys
TOKEN = "5787469608:AAGHPUMGPxUj0o8c0BYmEMXrtnkj6b6_lvw"
chat_id = "-1001890296042"
message = "hello from your telegram bot"

import urllib

with open('/tmp/data.txt', 'r') as file:
    data = file.read()




def send(text):

    safe_string = urllib.parse.quote_plus(text)

    url = f"https://api.telegram.org/bot{TOKEN}/sendMessage?chat_id={chat_id}&text={safe_string}&parse_mode=markdown"
    response = requests.get(url)


    if not response.ok:
        print(response.json())

send(data)