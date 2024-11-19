from fastapi import FastAPI


app = FastAPI()


@app.get("/get_things_list")
def get_things_list(type: str):
    if type == "lost":
        pass
    elif type == "found":
        pass
    

@app.post("/add_new_thing")
def add_new_thing(type: str):
    if type == "lost":
        pass
    elif type == "found":
        pass


@app.get("/change_thing_status")
def change_thing_status(id: int):
    pass
    
