import type { Component } from "solid-js";
import { createSignal } from "solid-js";

import styles from "./addnewlostthing.module.css";
import fileToBase64 from "./fileToBase64";

interface AddNewLostThingProps {
  onClick: func;
}

interface PostLostThingDataProps {
  thingName: string;
  userContacts: string;
  customText: string;
  thingPhoto: string;
}

interface PostFoundThingDataProps {
  thingName: string;
  thingLocation: string;
  customText: string;
  thingPhoto: string;
}

const postLostThingData = async (data: PostLostThingDataProps) => {
  const response = await fetch(`http://localhost:8000/add_new_lost_thing`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json; charset=utf-8",
    },
    body: JSON.stringify(data),
  });
  return response.json();
};

const postFoundThingData = async (data: PostFoundThingDataProps) => {
  const response = await fetch(`http://localhost:8000/add_new_found_thing`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json; charset=utf-8",
    },
    body: JSON.stringify(data),
  });
  return response.json();
};

const AddNewLostThing: Component = (props: AddNewLostThingProps) => {
  const [chooseThingType, setChooseThingType] = createSignal(true);
  const [addNewLostThing, setAddNewLostThing] = createSignal(false);
  const [addNewFoundThing, setAddNewFoundThing] = createSignal(false);

  const [thingName, setThingName] = createSignal("");
  const [userContacts, setUserContacts] = createSignal("");
  const [thingLocation, setThingLocation] = createSignal("");
  const [customText, setCustomText] = createSignal("");
  const [thingPhoto, setThingPhoto] = createSignal();

  const [data, setData] = createSignal({});

  return (
    <div class={styles.wrapper}>
      <div class={styles.box}>
        {chooseThingType() && (
          <>
            <button
              onClick={() => {
                setChooseThingType(false);
                setAddNewLostThing(true);
              }}
            >
              Я потерял вещь
            </button>
            <button
              onClick={() => {
                setChooseThingType(false);
                setAddNewFoundThing(true);
              }}
            >
              Я нашёл вещь
            </button>
          </>
        )}
        {addNewLostThing() && (
          <>
            <div class={styles.box_title}>Добавить потерянную вещь</div>
            <form method="post">
              <input
                placeholder="Что Вы потеряли?"
                value={thingName()}
                onInput={(e) => setThingName(e.target.value)}
                required
              />
              <input
                placeholder="Как с Вами можно связаться?"
                value={userContacts()}
                onInput={(e) => setUserContacts(e.target.value)}
                required
              />
              <textarea
                placeholder="Здесь можно оставить сообщение"
                value={customText()}
                onInput={(e) => setCustomText(e.target.value)}
                required
              />
              {thingPhoto() && <img src={thingPhoto()} />}
              <input
                type="file"
                accept="image/jpeg"
                onInput={(e) =>
                  fileToBase64(e.target.files[0]).then((r) => setThingPhoto(r))
                }
              />
              <button
                onClick={(e) => {
                  e.preventDefault();
                  setData({
                    thing_name: thingName(),
                    user_contacts: userContacts(),
                    custom_text: customText(),
                    thing_photo: thingPhoto(),
                  });
                  postLostThingData(data());
                  window.location.reload();
                }}
              >
                Отправить
              </button>
            </form>
          </>
        )}
        {addNewFoundThing() && (
          <>
            <div class={styles.box_title}>Добавить найденную вещь</div>
            <form method="post">
              <input
                placeholder="Что Вы нашли?"
                value={thingName()}
                onInput={(e) => setThingName(e.target.value)}
                required
              />
              <input
                placeholder="Где забрать вещь?"
                value={thingLocation()}
                onInput={(e) => setThingLocation(e.target.value)}
                required
              />
              <textarea
                placeholder="Здесь можно оставить сообщение"
                value={customText()}
                onInput={(e) => setCustomText(e.target.value)}
                required
              />
              {thingPhoto() && <img src={thingPhoto()} />}
              <input
                type="file"
                accept="image/jpeg"
                onInput={(e) =>
                  fileToBase64(e.target.files[0]).then((r) => setThingPhoto(r))
                }
              />
              <button
                onClick={(e) => {
                  e.preventDefault();
                  setData({
                    thing_name: thingName(),
                    thing_location: thingLocation(),
                    custom_text: customText(),
                    thing_photo: thingPhoto(),
                  });
                  postFoundThingData(data());
                  window.location.reload();
                }}
              >
                Отправить
              </button>
            </form>
          </>
        )}
      </div>
      <div
        class={styles.background}
        onClick={props.onClick}
      ></div>
    </div>
  );
};

export default AddNewLostThing;
