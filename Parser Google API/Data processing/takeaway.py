import pandas as pd
import scipy.spatial
import numpy as np



df = pd.read_csv("data_food3.csv")


df['opportunity_take_away'] = [1 if 'meal_takeaway' in x else 0 for x in df['types']]

df = df.drop(columns=['place_id'])


df.to_csv('data_food_final.csv', index=False, encoding='utf-8')
