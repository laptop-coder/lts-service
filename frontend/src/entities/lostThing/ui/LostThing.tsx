import type { Component } from "solid-js";
import { createSignal } from "solid-js";

import "../../../app/styles.css";
import { months } from "../../../shared/constants/index";
import { changeThingStatus } from "../api/changeThingStatus";
import { Props } from "../model/Props";

export const LostThing: Component = ({ syncList, tabIndex, props }: Props) => {
  const monthNumber = Number(props.publication_date.slice(5, 7));
  const day = Number(props.publication_date.slice(8, 10));
  const month = months[monthNumber - 1];
  const year = props.publication_date.slice(0, 4);
  const time = props.publication_time;
  const [thingHidden, setThingHidden] = createSignal(false);

  return (
    <div class={thingHidden() ? "thing thing__hidden" : "thing"}>
      <div class="thing__title">
        {props.thing_name} (№{props.id})
      </div>
      <div class="thing__content">
        <div>
          Опубликовано:{" "}
          <b>
            {day} {month} {year} в {time}
          </b>
          <br />
          Email: <b>{props.email}</b>
          <br />
          <i>{props.custom_text}</i>
        </div>
        {props.thing_photo && (
          <img
            class="thing__photo"
            src={"data:image/jpeg;base64," + props.thing_photo}
          />
        )}
      </div>
      <button
        tabindex={tabIndex}
        onClick={() => {
          setThingHidden(true);
          setTimeout(() => {
            changeThingStatus(props.id).then(() => syncList());
          }, 500);
        }}
      >
        Я забрал вещь
      </button>
    </div>
  );
};
