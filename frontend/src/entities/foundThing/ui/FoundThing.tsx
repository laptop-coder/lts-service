import type { Component } from "solid-js";

import "../../../app/styles.css";
import { months } from "../../../shared/constants/index";
import { changeThingStatus } from "../api/changeThingStatus";
import { Props } from "../types/Props";

export const FoundThing: Component = ({ tabIndex, props }: Props) => {
  const monthNumber = Number(props.publication_date.slice(5, 7));

  const day = Number(props.publication_date.slice(8, 10));
  const month = months[monthNumber - 1];
  const year = props.publication_date.slice(0, 4);

  const time = props.publication_time;

  return (
    <div class="thing">
      <div class="thing__title">
        {props.thing_name} (№{props.id})
      </div>
      <div class="thing__content">
        Опубликовано: {day} {month} {year} в {time}
        <br />
        Где забрать: {props.thing_location}
        <br />
        {props.custom_text}
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
          changeThingStatus(props.id);
          window.location.reload();
        }}
      >
        Я забрал вещь
      </button>
    </div>
  );
};
