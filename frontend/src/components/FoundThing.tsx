import type { Component } from "solid-js";

import "../styles.css";
import { months } from "../constants";

interface FoundThingProps {
  id: number;
  publication_date: string;
  publication_time: string;
  thing_name: string;
  thing_location: string;
  custom_text: string;
}

const handleClick = async (type: string, id: number) => {
  const response = await fetch(
    `http://localhost:8000/change_thing_status?type=${type}&id=${id}`,
  );
  return response.json();
};

const FoundThing: Component = (props: FoundThingProps) => {
  const monthNumber = Number(props.props.publication_date.slice(5, 7));

  const day = Number(props.props.publication_date.slice(8, 10));
  const month = months[monthNumber - 1];
  const year = props.props.publication_date.slice(0, 4);

  const time = props.props.publication_time;

  return (
    <div class="thing">
      <div class="thing__title">
        {props.props.thing_name} (№ {props.props.id})
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
          handleClick("found", props.props.id);
          window.location.reload();
        }}
      >
        Я забрал вещь
      </button>
    </div>
  );
};

export default FoundThing;
