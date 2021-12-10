import pandas as pd
import scipy.spatial
import numpy as np
from geopy.distance import vincenty
from geopy.distance import geodesic

df = pd.read_csv("data_food.csv")

df1 = pd.read_json("places_museum1.json")
df2 = pd.read_json("places_art_gallery.json")
df3 = pd.read_json("places_shopping_mall1.json")
df4 = pd.read_json("places_tourist_attraction1.json")

list_of_location = [df1, df2, df3, df4]
list_of_name = ["museum", "art_gallery", "shopping_mall", "places_tourist_attraction"]

for i in range(len(list_of_location)):
    list_of_location[i]['location'] = pd.DataFrame(list_of_location[i].geometry.tolist()).location
    list_of_location[i] = list_of_location[i].join(pd.DataFrame(list_of_location[i].location.tolist()))

    mat1 = scipy.spatial.distance.cdist(df[["lat", "lng"]], list_of_location[i][["lat", "lng"]], lambda u, v: vincenty(u, v).kilometers)

    new_df = pd.DataFrame(mat1, index=df["id"], columns=list_of_location[i]["id"])

    closest = np.where(new_df < 0.8, 1, 0)

    df.loc[:,f'Quantity of nearest {list_of_name[i]}'] = closest.sum(axis=1)
    # df.loc[f'Quantity of nearest {list_of_name[i]}', :] = closest.sum(axis=0)

df.to_csv('data_food2.csv', index=False, encoding='utf-8')



