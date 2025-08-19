import { JSX } from 'solid-js';

import styles from './ChangeThingStatusButton.module.css';
import changeThingStatus from '../../utils/changeThingStatus';
import type LostThing from '../../types/LostThing';
import type FoundThing from '../../types/FoundThing';
import type ThingType from '../../types/ThingType';

const ChangeThingStatusButton = (
  props: JSX.ButtonHTMLAttributes<HTMLButtonElement> & {
    thingType: ThingType;
    thingId: number;
    reloadLostThingsList: (
      info?: unknown,
    ) => LostThing[] | Promise<LostThing[] | undefined> | null | undefined;
    reloadFoundThingsList: (
      info?: unknown,
    ) => FoundThing[] | Promise<FoundThing[] | undefined> | null | undefined;
  },
): JSX.Element => (
  <button
    class={styles.change_thing_status_button}
    onclick={() => {
      changeThingStatus({ thingType: props.thingType, thingId: props.thingId });
      props.reloadLostThingsList();
      props.reloadFoundThingsList();
    }}
    type='button'
  >
    Я забрал вещь
  </button>
);

export default ChangeThingStatusButton;
