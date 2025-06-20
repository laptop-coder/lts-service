import type { Component } from 'solid-js';
import { createSignal } from 'solid-js';

import '../../../app/styles.css';
import { months } from '../../../shared/constants/index';
import { changeThingStatus } from '../api/changeThingStatus';
import { LostThingProps } from '../model/ThingProps';
import { FoundThingProps } from '../model/ThingProps';

export const Thing: Component<LostThingProps & FoundThingProps> = (props) => {
  const custom_text = props.custom_text;
  const email = props.email;
  const id = props.id;
  const page = props.page;
  const syncList = props.syncList;
  const tabIndex = props.tabIndex;
  const thing_location = props.thing_location;
  const thing_name = props.thing_name;
  const type = props.type;

  const monthNumber = Number(props.publication_date.slice(5, 7));
  const day = Number(props.publication_date.slice(8, 10));
  const month = months[monthNumber - 1];
  const year = props.publication_date.slice(0, 4);
  const time = props.publication_time;
  const [thingHidden, setThingHidden] = createSignal(false);

  const path_to_photo = `/storage/${type}/${id}.jpeg`;

  // Check if the photo exists and put the result in photoExists()
  const [photoExists, setPhotoExists] = createSignal(false);
  const photo = new Image();
  photo.onload = () => setPhotoExists(true);
  photo.src = path_to_photo;

  return (
    <div class={thingHidden() ? 'thing thing__hidden' : 'thing'}>
      <div class='thing__title'>
        {thing_name} (№{id})
      </div>
      <div class='thing__content'>
        <div>
          Опубликовано:{' '}
          <b>
            {day} {month} {year} в {time}
          </b>
          <br />
          {type === 'lost' && (
            <>
              Email: <b>{email}</b>
            </>
          )}
          {type === 'found' && (
            <>
              Где забрать: <b>{thing_location}</b>
            </>
          )}
          <br />
          <i>{custom_text}</i>
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
        tabIndex={tabIndex}
        onClick={() => {
          if (confirm('Вы уверены?')) {
            setThingHidden(true);
            setTimeout(() => {
              changeThingStatus(type, id).then(() => syncList());
            }, 500);
          }
        }}
      >
        Я забрал вещь
      </button>
    </div>
  );
};
