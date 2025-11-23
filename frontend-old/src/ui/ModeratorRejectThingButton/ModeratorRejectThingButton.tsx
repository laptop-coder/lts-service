import { JSX } from 'solid-js';

import styles from './ModeratorRejectThingButton.module.css';
import verifyThing from '../../utils/verifyThing';
import { ASSETS_ROUTE } from '../../utils/consts';
import { CONFIRM_ACTION_MESSAGE } from '../../utils/consts';
import type LostThing from '../../types/LostThing';
import type FoundThing from '../../types/FoundThing';
import type ThingType from '../../types/ThingType';

const ModeratorRejectThingButton = (
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
    class={styles.reject_button}
    onclick={() => {
      if (confirm(CONFIRM_ACTION_MESSAGE)) {
        verifyThing({
          thingType: props.thingType,
          thingId: props.thingId,
          action: 'reject',
        });
        props.reloadLostThingsList();
        props.reloadFoundThingsList();
      }
    }}
    type='button'
  >
    <img
      src={`${ASSETS_ROUTE}/reject.svg`}
      class={styles.img}
    />
  </button>
);

export default ModeratorRejectThingButton;
