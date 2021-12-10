import pandas as pd
import scipy.spatial
import numpy as np
from geopy.distance import vincenty
import json

df1 = pd.read_json("places_tourist_attraction.json")

df = df1[df1.types.map(lambda x: "bakery" in x)]


#df.to_json("places_museum1.json", orient='records')

#with open("places_tourist_attraction1.json", 'w') as f:
#    json.dump(a, f, indent=2)

#df = pd.read_json("places_museum1.json")
print(len(df))
print(df.name)

df2 = pd.read_json("places_bakery.json")
df2 = df2[df2.types.map(lambda x: "bakery" in x)]
print(len(df2))
print(df2.name)