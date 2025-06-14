import type { Component } from 'solid-js';
import { createSignal } from 'solid-js';

import '../../../app/styles.css';
import { months } from '../../../shared/constants/index';
import { changeThingStatus } from '../api/changeThingStatus';
import { LostThingProps } from '../model/LostThingProps';

export const ModeratorPageLostThing: Component<LostThingProps> = (props) => {
  const monthNumber = Number(props.publication_date.slice(5, 7));
  const day = Number(props.publication_date.slice(8, 10));
  const month = months[monthNumber - 1];
  const year = props.publication_date.slice(0, 4);
  const time = props.publication_time;
  const [thingHidden, setThingHidden] = createSignal(false);

  const path_to_photo = `/storage/lost/${props.id}.jpeg`;

  // Check if the photo exists and put the result in photoExists()
  const [photoExists, setPhotoExists] = createSignal(false);
  const photo = new Image();
  photo.onload = () => setPhotoExists(true);
  photo.src = path_to_photo;

  return (
    <div class={thingHidden() ? 'thing thing__hidden' : 'thing'}>
      <div class='thing__title'>
        {props.thing_name} (№{props.id})
      </div>
      <div class='thing__content'>
        <div>
          Опубликовано:{' '}
          <b>
            {day} {month} {year} в {time}
          </b>
          <br />
          Email: <b>{props.email}</b>
          <br />
          <i>{props.custom_text}</i>
        </div>
        {photoExists() && (
          <img
            class='thing__photo'
            src={path_to_photo}
            onClick={(event) => event.target.requestFullscreen()}
          />
        )}
      </div>
      <button
        tabIndex={props.tabIndex}
        onClick={() => {
          if (confirm('Вы уверены?')) {
            setThingHidden(true);
            setTimeout(() => {
              changeThingStatus(props.id).then(() => props.syncList());
            }, 500);
          }
        }}
      >
        Я забрал вещь
      </button>
    </div>
  );
};
