import os
import json
import requests
import pics_pb2
from io import BytesIO
from PIL import Image

host = "http://localhost:8081" or os.environ["SERVICE_URL"]
url = host + "/api/v1/pictures"

def getDataFromDB(type_:str, number:int):
    resp = requests.get(url, params={"type": type_, "number": number})
    buffer = resp.content
    pics = pics_pb2.Pics()
    pics.ParseFromString(buffer)
    # 此处buffer无法释放，会产生内存泄漏，仅供demo使用
    return [(pic.id, pic.location, Image.open(BytesIO(pic.pic_data))) for pic in pics.pic]


def sendDataToDB(key:int, is_valid:bool):
    data = {
        "id": key,
        "is_valid": is_valid
    }
    resp = requests.post(url, json=data)
    print(json.loads(resp.text))