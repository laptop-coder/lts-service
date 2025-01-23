import type { Component } from "solid-js";
import { createSignal } from "solid-js";

import "../app/styles.css";
import { fileToBase64 } from "../shared/lib/utils/index";
import { POST } from "../shared/lib/utils/index";
import { d } from "../shared/assets/index";
import { SVG } from "../shared/ui/index";

interface LostThingData {
  thingName: string;
  userContacts: string;
  customText: string;
  thingPhoto: string;
}

interface FoundThingData {
  thingName: string;
  thingLocation: string;
  customText: string;
  thingPhoto: string;
}

const checkLostThingDataType = (data: LostThingData) => {
  return true;
};

const checkFoundThingDataType = (data: FoundThingData) => {
  return true;
};

const AddNewThing: Component = () => {
  const [chooseThingType, setChooseThingType] = createSignal(true);
  const [addNewLostThing, setAddNewLostThing] = createSignal(false);
  const [addNewFoundThing, setAddNewFoundThing] = createSignal(false);

  const [thingName, setThingName] = createSignal("");
  const [userContacts, setUserContacts] = createSignal("");
  const [thingLocation, setThingLocation] = createSignal("");
  const [customText, setCustomText] = createSignal("");
  const [thingPhoto, setThingPhoto] = createSignal();

  const [data, setData] = createSignal({});

  const [uploadPhotoFocus, setUploadPhotoFocus] = createSignal(false);

  return (
    <>
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
          <div class="box_title">Добавить потерянную вещь</div>
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
            {thingPhoto() && (
              <img
                class="thing__photo"
                src={thingPhoto()}
              />
            )}
            <input
              type="file"
              class="hidden"
              id="upload-photo__input"
              accept="image/jpeg"
              onFocus={() => setUploadPhotoFocus((prev) => !prev)}
              onBlur={() => setUploadPhotoFocus((prev) => !prev)}
              onInput={(e) =>
                fileToBase64(e.target.files[0]).then((r) => setThingPhoto(r))
              }
            />
            <label
              class={`upload-photo__label${thingPhoto() ? " hidden" : ""}${uploadPhotoFocus() ? " focus" : ""}`}
              for="upload-photo__input"
            >
              Выберите файл
            </label>
            <button
              onClick={(e) => {
                e.preventDefault();
                setData({
                  thing_name: thingName(),
                  user_contacts: userContacts(),
                  custom_text: customText(),
                  thing_photo: thingPhoto(),
                });
                if (checkLostThingDataType()) {
                  POST("add_new_lost_thing", data());
                } else {
                  console.log("Type error (POST, lost things)");
                }
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
          <div class="box_title">Добавить найденную вещь</div>
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
            {thingPhoto() && (
              <img
                class="thing__photo"
                src={thingPhoto()}
              />
            )}
            <input
              type="file"
              class="hidden"
              id="upload-photo__input"
              accept="image/jpeg"
              onFocus={() => setUploadPhotoFocus((prev) => !prev)}
              onBlur={() => setUploadPhotoFocus((prev) => !prev)}
              onInput={(e) =>
                fileToBase64(e.target.files[0]).then((r) => setThingPhoto(r))
              }
            />
            <label
              class={`upload-photo__label${thingPhoto() ? " hidden" : ""}${uploadPhotoFocus() ? " focus" : ""}`}
              for="upload-photo__input"
            >
              Выберите файл
            </label>
            <button
              onClick={(e) => {
                e.preventDefault();
                setData({
                  thing_name: thingName(),
                  thing_location: thingLocation(),
                  custom_text: customText(),
                  thing_photo: thingPhoto(),
                });
                if (checkFoundThingDataType()) {
                  POST("add_new_found_thing", data());
                } else {
                  console.log("Type error (POST, found things)");
                }
                window.location.reload();
              }}
            >
              Отправить
            </button>
          </form>
        </>
      )}
    </>
  );
};

export default AddNewThing;
