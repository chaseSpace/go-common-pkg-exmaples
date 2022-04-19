import requests, fake_useragent


def x():
    ua = fake_useragent.UserAgent()
    proxies = {
        "https": "socks5://root:12345678@43.156.75.168:8008"
    }
    headers = {
        "user-agent": ua.firefox
    }
    res = requests.get("https://checkip.amazonaws.com/", proxies=proxies, headers=headers, timeout=3)
    print(res.text)


if __name__ == '__main__':
    x()
