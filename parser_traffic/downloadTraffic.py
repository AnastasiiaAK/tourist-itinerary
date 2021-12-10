from google.transit import gtfs_realtime_pb2
import requests
import datetime
import time
from time import gmtime, strftime
from datetime import date
import json
import os
from google.protobuf.json_format import MessageToDict



dat = date.today().strftime('%d')
fileData = date.today().strftime('%d%h')
data = date.today().strftime('%d')

def traffic(dat, fileData, data):
    while int(data) < int(dat) + 14:
        
        feed = gtfs_realtime_pb2.FeedMessage()
        response = requests.get('http://transport.orgp.spb.ru/Portal/transport/internalapi/gtfs/realtime/vehicle')
        for entity in feed.entity:
            if entity.HasField('trip_update'):
                print(entity.trip_update)
        feed.ParseFromString(response.content)
        bus = feed.entity
        messageToDict = MessageToDict(feed, preserving_proto_field_name = True)



        times = datetime.datetime.now().strftime('%H:%M')
        weekday = date.today().strftime('%A')
        data = date.today().strftime('%d')

        if int(data) > int(dat) + 6:
            fileData = date.today().strftime('%d%h')


        directory = f'/Users/user/Desktop/traffic_road/{fileData}/{weekday}'

        if not os.path.exists(directory):
            os.makedirs(directory)    


        with open(f'/Users/user/Desktop/traffic_road/{fileData}/{weekday}/{times}.txt', 'w') as file:
             file.write(json.dumps(messageToDict))
                
        
        time.sleep(180)


try:
    traffic(dat, fileData, data)
except:
    traffic(dat, fileData, data)
