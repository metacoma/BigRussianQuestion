import vk
import requests 
import sys
import os

album_id = int(os.getenv('VK_ALBUM_ID'))
group_id = int(os.getenv('VK_GROUP_ID'))
token = os.getenv('VK_TOKEN')
filename = sys.argv[1]
api = vk.API(vk.Session(access_token=token), v=5.122,scope=+4)
upload_url = api.photos.getWallUploadServer(group_id=group_id)['upload_url'] 
resp = requests.post(upload_url, files = {'file': open(filename, 'rb')}).json()
s = api.photos.saveWallPhoto(group_id=group_id, server = resp['server'], photo= resp['photo'], hash = resp['hash'])
api.wall.post(owner_id = -group_id, from_group = 1, message=" ", attachments=f"photo{s[0]['owner_id']}_{s[0]['id']}")
