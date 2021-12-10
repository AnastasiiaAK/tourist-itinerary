import pandas as pd
import numpy as np
from geopy.distance import geodesic
import scipy.spatial



df = pd.read_csv("data_food2.csv")
#print(df.head())

#print(df.types)
print(len(df))

df1 = df[df['types'].apply(lambda x: 'cafe' in x)]

df2 = df[df['types'].apply(lambda x: 'restaurant' in x and 'cafe' not in x )]

df3 = df[df['types'].apply(lambda x: 'bar' in x and 'restaurant' not in x and 'cafe' not in x)]

list_of_location = [df1, df2, df3]

for i in range(len(list_of_location)):

    mat1 = scipy.spatial.distance.cdist(list_of_location[i][["lat", "lng"]], list_of_location[i][["lat", "lng"]], lambda u, v: geodesic(u, v).kilometers)

    new_df = pd.DataFrame(mat1, index= list_of_location[i]["id"], columns=list_of_location[i]["id"])

    closest = np.where(new_df < 0.8, 1, 0)

    list_of_location[i].loc[:, f'Quantity of nearest same establishment'] = closest.sum(axis=1)


df = pd.concat([df1, df2, df3], ignore_index=True)

df.to_csv('data_food3.csv', index=False, encoding='utf-8')


