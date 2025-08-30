import { JSX } from 'solid-js';

import styles from './ModeratorAcceptThingButton.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';
import verifyThing from '../../utils/verifyThing';
import { CONFIRM_ACTION_MESSAGE } from '../../utils/consts';
import type LostThing from '../../types/LostThing';
import type FoundThing from '../../types/FoundThing';
import type ThingType from '../../types/ThingType';

const ModeratorAcceptThingButton = (
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
    class={styles.accept_button}
    onclick={() => {
      if (confirm(CONFIRM_ACTION_MESSAGE)) {
        verifyThing({
          thingType: props.thingType,
          thingId: props.thingId,
          action: 'accept',
        });
        props.reloadLostThingsList();
        props.reloadFoundThingsList();
      }
    }}
    type='button'
  >
    <img
      src={`${ASSETS_ROUTE}/accept.svg`}
      class={styles.img}
    />
  </button>
);

export default ModeratorAcceptThingButton;
