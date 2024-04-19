import hashlib
import requests

input = input()
r = requests.get("https://stepik.org/media/attachments/lesson/668860/dictionary.txt")
dict = r.text.lower().split("\n")


for word in dict:
    word = word.strip()
    hf = hashlib.md5()
    hf.update(word.encode())
    h = hf.hexdigest()
    
    if h == input:
        print(word)
        break
