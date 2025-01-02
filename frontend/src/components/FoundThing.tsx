import type { Component } from "solid-js";

import Button from "./Button";
import button_styles from "./button.module.css";
import styles from "./thing.module.css";
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
    <div class={styles.thing}>
      <div class={styles.thing__title}>
        {props.props.id}. {props.props.thing_name}
      </div>
      <div class={styles.thing__content}>
        Опубликовано: {day} {month} {year} в {time}
        <br />
        Где забрать: {props.props.thing_location}
        <br />
        {props.props.custom_text}
        {props.props.thing_photo && (
          <img
            class={styles.thing__photo}
            src={"data:image/jpeg;base64," + props.props.thing_photo}
          />
        )}
      </div>
      <Button
        class={button_styles.wide_button}
        onClick={() => handleClick("found", props.props.id)}
        type="text"
        value="Я забрал вещь"
      />
    </div>
  );
};

export default FoundThing;
