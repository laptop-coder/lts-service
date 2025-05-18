import type { Component, Setter, Accessor } from "solid-js";
import { createSignal } from "solid-js";
import { autofocus } from "@solid-primitives/autofocus";

import type { LostThingData } from "../model/LostThingData";
import type { FoundThingData } from "../model/FoundThingData";
import "../../../app/styles.css";
import { fileToBase64 } from "../../../shared/lib/utils/index";
import { POST } from "../../../shared/lib/utils/index";
import { d } from "../../../shared/assets/index";
import { SVG } from "../../../shared/ui/index";
import { Props } from "../model/Props";

const handleLostThingButtonClick = (
  setChooseThingType: Setter<boolean>,
  setAddNewLostThing: Setter<boolean>,
) => {
  setChooseThingType(false);
  setAddNewLostThing(true);
};

const handleFoundThingButtonClick = (
  setChooseThingType: Setter<boolean>,
  setAddNewFoundThing: Setter<boolean>,
) => {
  setChooseThingType(false);
  setAddNewFoundThing(true);
};

const keyDown = (
  event,
  chooseThingType: Accessor<boolean>,
  setChooseThingType: Setter<boolean>,
  setAddNewLostThing: Setter<boolean>,
  setAddNewFoundThing: Setter<boolean>,
) => {
  switch (event.key) {
    case "l":
      if (chooseThingType())
        handleLostThingButtonClick(setChooseThingType, setAddNewLostThing);
      break;
    case "f":
      if (chooseThingType())
        handleFoundThingButtonClick(setChooseThingType, setAddNewFoundThing);
      break;
  }
};

const checkLostThingDataType = (data: LostThingData) => {
  return true;
};

const checkFoundThingDataType = (data: FoundThingData) => {
  return true;
};

export const AddNewThing: Component = ({
  syncLostThingsList,
  syncFoundThingsList,
  setAddNewThing,
}: Props) => {
  const [chooseThingType, setChooseThingType] = createSignal(true);
  const [addNewLostThing, setAddNewLostThing] = createSignal(false);
  const [addNewFoundThing, setAddNewFoundThing] = createSignal(false);

  const [thingName, setThingName] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [thingLocation, setThingLocation] = createSignal("");
  const [customText, setCustomText] = createSignal("");
  const [thingPhoto, setThingPhoto] = createSignal();

  const [data, setData] = createSignal({});

  const [uploadPhotoFocus, setUploadPhotoFocus] = createSignal(false);

  return (
    <>
      {chooseThingType() && (
        <div
          class="choose_thing_type"
          tabindex="1"
          autofocus
          use:autofocus
          onKeyDown={(event) =>
            keyDown(
              event,
              chooseThingType,
              setChooseThingType,
              setAddNewLostThing,
              setAddNewFoundThing,
            )
          }
        >
          <button
            onClick={() =>
              handleLostThingButtonClick(setChooseThingType, setAddNewLostThing)
            }
          >
            Я потерял вещь
          </button>
          <button
            onClick={() =>
              handleFoundThingButtonClick(
                setChooseThingType,
                setAddNewFoundThing,
              )
            }
          >
            Я нашёл вещь
          </button>
        </div>
      )}
      {addNewLostThing() && (
        <>
          <div class="box_title">Добавить потерянную вещь</div>
          <form method="post">
            <input
              placeholder="Что Вы потеряли?*"
              value={thingName()}
              onInput={(event) => setThingName(event.target.value)}
              autofocus
              use:autofocus
            />
            <input
              type="email"
              placeholder="Email*"
              value={email()}
              onInput={(event) => setEmail(event.target.value)}
            />
            <textarea
              placeholder="Здесь можно оставить сообщение"
              value={customText()}
              onInput={(event) => setCustomText(event.target.value)}
            />
            {thingPhoto() && (
              <img
                class="thing__photo"
                src={thingPhoto()}
                onClick={(event) => event.target.requestFullscreen()}
              />
            )}
            <input
              type="file"
              class="hidden"
              id="upload-photo__input"
              accept="image/jpeg"
              onFocus={() => setUploadPhotoFocus((prev) => !prev)}
              onBlur={() => setUploadPhotoFocus((prev) => !prev)}
              tabindex={thingPhoto() ? "-1" : "0"}
              onInput={(event) =>
                fileToBase64(event.target.files[0]).then((photoBase64) =>
                  setThingPhoto(photoBase64),
                )
              }
            />
            <label
              class={`upload-photo__label${thingPhoto() ? " hidden" : ""}${uploadPhotoFocus() ? " focus" : ""}`}
              for="upload-photo__input"
            >
              Выберите файл
            </label>
            <button
              onClick={(event) => {
                event.preventDefault();
                if (thingName() !== "" && email() !== "") {
                  setData({
                    thing_name: thingName(),
                    email: email(),
                    custom_text: customText(),
                    thing_photo: thingPhoto(),
                  });
                  if (checkLostThingDataType()) {
                    POST("add_new_lost_thing", data()).then(() =>
                      syncLostThingsList(),
                    );
                  } else {
                    console.log("Type error (POST, lost things)");
                  }
                  setAddNewThing(false);
                } else {
                  alert("Обязательные поля не заполнены");
                }
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
              placeholder="Что Вы нашли?*"
              value={thingName()}
              onInput={(event) => setThingName(event.target.value)}
              autofocus
              use:autofocus
            />
            <input
              placeholder="Где забрать вещь?*"
              value={thingLocation()}
              onInput={(event) => setThingLocation(event.target.value)}
            />
            <textarea
              placeholder="Здесь можно оставить сообщение"
              value={customText()}
              onInput={(event) => setCustomText(event.target.value)}
            />
            {thingPhoto() && (
              <img
                class="thing__photo"
                src={thingPhoto()}
                onClick={(event) => event.target.requestFullscreen()}
              />
            )}
            <input
              type="file"
              class="hidden"
              id="upload-photo__input"
              accept="image/jpeg"
              onFocus={() => setUploadPhotoFocus((prev) => !prev)}
              onBlur={() => setUploadPhotoFocus((prev) => !prev)}
              tabindex={thingPhoto() ? "-1" : "0"}
              onInput={(event) =>
                fileToBase64(event.target.files[0]).then((photoBase64) =>
                  setThingPhoto(photoBase64),
                )
              }
            />
            <label
              class={`upload-photo__label${thingPhoto() ? " hidden" : ""}${uploadPhotoFocus() ? " focus" : ""}`}
              for="upload-photo__input"
            >
              Выберите файл
            </label>
            <button
              onClick={(event) => {
                event.preventDefault();
                if (thingName() !== "" && thingLocation() !== "") {
                  setData({
                    thing_name: thingName(),
                    thing_location: thingLocation(),
                    custom_text: customText(),
                    thing_photo: thingPhoto(),
                  });
                  if (checkFoundThingDataType()) {
                    POST("add_new_found_thing", data()).then(() =>
                      syncFoundThingsList(),
                    );
                  } else {
                    console.log("Type error (POST, found things)");
                  }
                  setAddNewThing(false);
                } else {
                  alert("Обязательные поля не заполнены");
                }
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
