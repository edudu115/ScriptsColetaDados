#!/usr/bin/python3
from pymongo import MongoClient # type: ignore
import sys
import re
import json

ip=sys.argv[1]
user=sys.argv[2]
passw=sys.argv[3]
dbInfo=sys.argv[4]
collectionInfo=sys.argv[5]
indexInfo=sys.argv[6]

data = {"data": []}

try:
    client = MongoClient("mongodb://"+ip+":27017/",
                    username=user,
                    password=passw)

    dbs = client.list_database_names()
    for db in dbs:
        if re.fullmatch(dbInfo, db):
            collections = client[db].list_collection_names()
            for collection in collections:
                if re.fullmatch(collectionInfo , collection):
                    index = client[db][collection].index_information()
                    for indexName, indexData in index.items():
                        if re.fullmatch(indexInfo, indexName):
                            status = 0 #Status 0 quer dizer alerta
                            if "expireAfterSeconds" in indexData:
                                if indexData["expireAfterSeconds"] > 0:
                                    status = 1 #Status 1 quer dizer ok
                                else:
                                    status = 0 #Status 0 quer dizer alerta
                            data["data"].append({
                                    "db": db,
                                    "collection": collection,
                                    "index": indexName,
                                    "status": status
                                })
    print(json.dumps(data, indent=4))
    client.close()
except Exception as e:
    print("Error: ", e)