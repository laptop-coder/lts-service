import { JSX } from 'solid-js';

import styles from './ThingDeleteButton.module.css';
import { ASSETS_ROUTE, Role, ThingType } from '../../utils/consts';
import deleteThing from '../../utils/deleteThing';

const ThingDeleteButton = (props: {
  thingName: string;
  thingType: ThingType;
  thingId: string;
  role: Role;
}): JSX.Element => (
  <button
    class={styles.thing_delete_button}
    onclick={() => {
      if (
        confirm(
          `Подтвердите удаление объявления "${props.thingName}". Это действие необратимо`,
        )
      ) {
        deleteThing({
          thingType: props.thingType,
          thingId: props.thingId,
          role: props.role,
        });
      }
    }}
    title='Удалить объявление'
  >
    <img src={`${ASSETS_ROUTE}/delete.svg`} />
  </button>
);

export default ThingDeleteButton;
