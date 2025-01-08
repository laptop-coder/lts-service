import type { Component } from "solid-js";

import "../../../app/styles.css";
import { months } from "../../../utils/constants";
import { changeThingStatus } from "../api/changeThingStatus";
import { Props } from "../types/Props";

export const FoundThing: Component = (props: Props) => {
  const monthNumber = Number(props.props.publication_date.slice(5, 7));

  const day = Number(props.props.publication_date.slice(8, 10));
  const month = months[monthNumber - 1];
  const year = props.props.publication_date.slice(0, 4);

  const time = props.props.publication_time;

  return (
    <div class="thing">
      <div class="thing__title">
        {props.props.thing_name} (№{props.props.id})
      </div>
      <div class="thing__content">
        Опубликовано: {day} {month} {year} в {time}
        <br />
        Где забрать: {props.props.thing_location}
        <br />
        {props.props.custom_text}
        {props.props.thing_photo && (
          <img
            class="thing__photo"
            src={"data:image/jpeg;base64," + props.props.thing_photo}
          />
        )}
      </div>
      <button
        onClick={() => {
          changeThingStatus(props.props.id);
          window.location.reload();
        }}
      >
        Я забрал вещь
      </button>
    </div>
  );
};
