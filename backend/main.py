from fastapi import FastAPI
import sqlite3

import config


app = FastAPI()

@app.get("/get_things_list")
def get_things_list(type: str):
    if type == "lost":
        connection = sqlite3.connect(config.PATH_TO_DB)
        with connection:
            cursor = connection.cursor()
            data = cursor.execute(
                """
                SELECT * FROM lost_thing WHERE status=0;
                """
            ).fetchall()
            formatted_data = []
            for i in range(len(data)):
                formatted_data.append({
                    "id": data[i][0],
                    "publication_date": data[i][1],
                    "publication_time": data[i][2],
                    "thing_name": data[i][3],
                    "user_contacts": data[i][4],
                    "path_to_thing_photo": data[i][5],
                    "custom_text": data[i][6],
                    "status": data[i][7]
                    })
            return formatted_data
    elif type == "found":
        connection = sqlite3.connect(config.PATH_TO_DB)
        with connection:
            cursor = connection.cursor()
            data = cursor.execute(
                """
                SELECT * FROM found_thing WHERE status=0;
                """
            ).fetchall()
            formatted_data = []
            for i in range(len(data)):
                formatted_data.append({
                    "id": data[i][0],
                    "publication_date": data[i][1],
                    "publication_time": data[i][2],
                    "thing_name": data[i][3],
                    "thing_location": data[i][4],
                    "path_to_thing_photo": data[i][5],
                    "custom_text": data[i][6],
                    "status": data[i][7]
                    })
            return formatted_data
    

@app.post("/add_new_thing")
def add_new_thing(type: str):
    if type == "lost":
        pass
    elif type == "found":
        pass


@app.get("/change_thing_status")
def change_thing_status(type: str, id: int):
    connection = sqlite3.connect(config.PATH_TO_DB)
    with connection:
        cursor = connection.cursor()
        cursor.execute(f
            """
            UPDATE {type}_thing SET status=1 WHERE id={id};
            """
        )

