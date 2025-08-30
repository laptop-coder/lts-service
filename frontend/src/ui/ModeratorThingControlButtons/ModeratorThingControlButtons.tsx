import { JSX } from 'solid-js';

import styles from './ModeratorThingControlButtons.module.css';
import changeThingStatus from '../../utils/changeThingStatus';
import { CONFIRM_ACTION_MESSAGE } from '../../utils/consts';
import { ASSETS_ROUTE } from '../../utils/consts';
import type LostThing from '../../types/LostThing';
import type FoundThing from '../../types/FoundThing';
import type ThingType from '../../types/ThingType';
import ModeratorAcceptThingButton from '../ModeratorAcceptThingButton/ModeratorAcceptThingButton';
import ModeratorRejectThingButton from '../ModeratorRejectThingButton/ModeratorRejectThingButton';

const ModeratorThingControlButtons = (
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
  <div class={styles.moderator_thing_control_buttons}>
    <ModeratorRejectThingButton
      thingType={props.thingType}
      thingId={props.thingId}
      reloadLostThingsList={props.reloadLostThingsList}
      reloadFoundThingsList={props.reloadFoundThingsList}
    />
    <ModeratorAcceptThingButton
      thingType={props.thingType}
      thingId={props.thingId}
      reloadLostThingsList={props.reloadLostThingsList}
      reloadFoundThingsList={props.reloadFoundThingsList}
    />
  </div>
);

export default ModeratorThingControlButtons;
