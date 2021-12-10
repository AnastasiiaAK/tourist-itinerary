import os
import json
import pandas as pd
df = pd.DataFrame()
your_path = '/Users/user/Desktop/traffic_road/22Oct/Thursday/'
files = os.listdir(your_path)
for file in files:
    with open(f'/Users/user/Desktop/traffic_road/22Oct/Thursday/{file}', 'r') as file:
        data = json.load(file)
    if len(data)>1:
        df0 = pd.io.json.json_normalize(data['entity'])
        df = df.append(df0, ignore_index=True)

print(df)

# опеределяем временной проежуток в 10 минут и фильтурем. с 08 00 до 08 10

import datetime
from datetime import datetime

# ниже закомментирован скрипт далее ля цикла по файлам
"""
#timestamp = datetime.datetime.fromtimestamp(1603281627)
#beginTime = timestamp.strftime('%H%M')
for hour in range(0, 23):
    for minute in range(0, 60, 10):
        currentIntervalBegin = f'2020-10-{day} {hour}:{minute}:00'
        currentIntervalEnd = f'2020-10-{day} {hour}:{minute+10}:00'

"""

# берем конкретный промежуток в 10 минкт
currentIntervalBegin = '2020-10-22 08:00:00'
currentIntervalEnd = '2020-10-22 08:10:00'
format_dateBegin = datetime.strptime(currentIntervalBegin, '%Y-%m-%d %H:%M:%S')
format_dateEnd = datetime.strptime(currentIntervalEnd, '%Y-%m-%d %H:%M:%S')
current_beginInterval = int(format_dateBegin.timestamp())
current_endInterval = int(format_dateEnd.timestamp())

# фильруем данные в нужном временном промежутке
currentDf0 = df[(df['vehicle.timestamp'].astype(int) > int(current_beginInterval))]
currentDf = currentDf0[currentDf0['vehicle.timestamp'].astype(int) < int(current_endInterval)]




currentDf['coordinates'] = currentDf['vehicle.position.latitude'].astype(str) + ','+ currentDf['vehicle.position.longitude'].astype(str)
allBus = pd.read_csv('~/public_transport.csv')
currentDf = currentDf.drop_duplicates()
listOfId = list(currentDf["id"].unique())

import geopy.distance

# создаем пустую таблицу
newDf = pd.DataFrame()
# для всех id из списка уникальных  id
for i in listOfId:
    # берем все данные для текущего id
    curId = currentDf[currentDf["id"] == i]
    route_id = curId["vehicle.trip.route_id"].iloc[0]
    busCur = allBus[allBus['route_id'] == int(route_id)]
    curId["nearest_stop"] = None
    curId["nearest_stop_distance"] = None

    # определяем для каждой точки ближайшую останвоку и расстояние до нее
    for row1 in range(len(curId)):
        minDist = 1000000
        for row2 in range(len(busCur)):
            if minDist > geopy.distance.geodesic(curId['coordinates'].iloc[row1], busCur['coordinates'].iloc[row2]).m:
                nearestStop = busCur['stop_id'].iloc[row2]
                minDist = geopy.distance.geodesic(curId['coordinates'].iloc[row1], busCur['coordinates'].iloc[row2]).m
            # расстояние до ближайшей остановки должно быть мень 400 метров, если это так, то добавляем эти данные в столбец о ближайшей остановке
            if minDist < 400:
                curId["nearest_stop"].iloc[row1] = nearestStop
                curId["nearest_stop_distance"].iloc[row1] = minDist

    # можно удалить эти данные из списка
    curId = curId[curId["nearest_stop"].notnull()]
    # переименовываем столбцы в более понятные названия
    curId = curId.rename(columns={"vehicle.trip.route_id": "route_id", "nearest_stop": "stop_id"})
    curId['route_id'] = curId['route_id'].astype(int)
    # удаляем дубликаты
    curId = curId.drop_duplicates(subset=['id', 'route_id', 'stop_id'], keep='first')

# если длина полученных данных меньше 1, то не имеет смысла искать разницу между предущие и следующей останвкой
# делаем merge по полученным данным от автобусов в настоящий момент времени и всеми данными автобусов

import geopy.distance

# создаем пустую таблицу
newDf = pd.DataFrame()
# для всех id из списка уникальных  id
for i in listOfId:
    # берем все данные для текущего id
    curId = currentDf[currentDf["id"] == i]
    right = curId[['stop_id', 'route_id', 'vehicle.timestamp', "nearest_stop_distance"]]

    left = busCur
    result = pd.merge(left, right, how='left', on=['stop_id', 'route_id'])
    # получаем остановки по порядку и данные о времени некоторых останвовок

    # бывает, что попадают данные из другого направления. Делаем, чтобы данные были толкьодля одного направления
    resultDirection = result[result["vehicle.timestamp"].notnull()]
    reverse = len(resultDirection[resultDirection["direction"] == "Обратное"])
    reverseIndex = resultDirection[resultDirection["direction"] == "Обратное"].index
    direct = len(resultDirection[resultDirection["direction"] == "Прямое"])
    directIndex = resultDirection[resultDirection["direction"] == "Прямое"].index
    if direct > reverse:
        for index in reverseIndex:
            result["vehicle.timestamp"].iloc[index] = None
    else:
        for index in directIndex:
            result["vehicle.timestamp"].iloc[index] = None

    # смотрим индексы для ненулевых данных
    indexes = result[result['vehicle.timestamp'].notnull()].index
    # делаем срез между ненулевыми данными
    resultCurrent = result.loc[min(indexes):max(indexes)]

    # считаем время между известными остановками
    commonTime = abs(int(resultCurrent["vehicle.timestamp"].loc[max(indexes)]) - int(
        resultCurrent["vehicle.timestamp"].loc[min(indexes)]))  # в секундах
    # считаем расстояние между остановками (они известы из фала о данных марурта автобусов)

    commonDist = resultCurrent.loc[min(indexes):(max(indexes) - 1)]["stop_distance"].sum()  # в км
    commonDistM = commonDist * 1000
    # считаем скорость

    speed = commonDistM / commonTime
    # заполняем данные о времени между остановками
    resultCurrent["timeBetweenStop"] = resultCurrent["stop_distance"] * 1000 / speed

    # добавляем в таблицу
    newDf = newDf.append(resultCurrent, ignore_index=True)
newDf = newDf[newDf["timeBetweenStop"].notnull()]
print(newDf)




import geopy.distance

# создаем пустую таблицу
newDf1 = pd.DataFrame()
# для всех id из списка уникальных  id
for i in listOfId:
    # берем все данные для текущего id
    curId = currentDf[currentDf["id"] == i]
    route_id = curId["vehicle.trip.route_id"].iloc[0]
    busCur = allBus[allBus['route_id'] == int(route_id)]
    curId["nearest_stop"] = None
    curId["nearest_stop_distance"] = None

    # определяем для каждой точки ближайшую останвоку и расстояние до нее
    for row1 in range(len(curId)):
        minDist = 1000000
        for row2 in range(len(busCur)):
            if minDist > geopy.distance.geodesic(curId['coordinates'].iloc[row1], busCur['coordinates'].iloc[row2]).m:
                nearestStop = busCur['stop_id'].iloc[row2]
                minDist = geopy.distance.geodesic(curId['coordinates'].iloc[row1], busCur['coordinates'].iloc[row2]).m
            # расстояние до ближайшей остановки должно быть мень 400 метров, если это так, то добавляем эти данные в столбец о ближайшей остановке
            if minDist < 400:
                curId["nearest_stop"].iloc[row1] = nearestStop
                curId["nearest_stop_distance"].iloc[row1] = minDist

    # можно удалить эти данные из списка
    curId = curId[curId["nearest_stop"].notnull()]
    # переименовываем столбцы в более понятные названия
    curId = curId.rename(columns={"vehicle.trip.route_id": "route_id", "nearest_stop": "stop_id"})
    curId['route_id'] = curId['route_id'].astype(int)
    # удаляем дубликаты
    curId = curId.drop_duplicates(subset=['id', 'route_id', 'stop_id'], keep='first')

    # если длина полученных данных меньше 1, то не имеет смысла искать разницу между предущие и следующей останвкой
    if len(curId) > 1:

        # делаем merge по полученным данным от автобусов в настоящий момент времени и всеми данными автобусов
        right = curId[['stop_id', 'route_id', 'vehicle.timestamp', "nearest_stop_distance"]]

        left = busCur
        result = pd.merge(left, right, how='left', on=['stop_id', 'route_id'])
        # получаем остановки по порядку и данные о времени некоторых останвовок

        # бывает, что попадают данные из другого направления. Делаем, чтобы данные были толкьодля одного направления
        resultDirection = result[result["vehicle.timestamp"].notnull()]
        reverse = len(resultDirection[resultDirection["direction"] == "Обратное"])
        reverseIndex = resultDirection[resultDirection["direction"] == "Обратное"].index
        direct = len(resultDirection[resultDirection["direction"] == "Прямое"])
        directIndex = resultDirection[resultDirection["direction"] == "Прямое"].index
        if direct > reverse:
            for index in reverseIndex:
                result["vehicle.timestamp"].iloc[index] = None
        else:
            for index in directIndex:
                result["vehicle.timestamp"].iloc[index] = None

        # смотрим индексы для ненулевых данных
        indexes = result[result['vehicle.timestamp'].notnull()].index
        # делаем срез между ненулевыми данными
        resultCurrent = result.loc[min(indexes):max(indexes)]

        # считаем время между известными остановками
        commonTime = abs(int(resultCurrent["vehicle.timestamp"].loc[max(indexes)]) - int(
            resultCurrent["vehicle.timestamp"].loc[min(indexes)]))
        # считаем расстояние между остановками (они известы из фала о данных марурта автобусов)

        commonDist = resultCurrent.loc[min(indexes):(max(indexes) - 1)]["stop_distance"].sum()

        # считаем скорость
        speed = commonDist / commonTime
        # заполняем данные о времени между остановками
        resultCurrent["timeBetweenStop"] = resultCurrent["stop_distance"] / speed

        # добавляем в таблицу
        newDf1 = newDf1.append(resultCurrent, ignore_index=True)
newDf1 = newDf1[newDf1["timeBetweenStop"].notnull()]
print(newDf1)