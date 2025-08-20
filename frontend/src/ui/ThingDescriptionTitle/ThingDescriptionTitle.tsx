import { JSX, ParentProps } from 'solid-js';

import { A } from '@solidjs/router';

import styles from './ThingDescriptionTitle.module.css';
import { THING_STATUS_ROUTE } from '../../utils/consts';
import type ThingType from '../../types/ThingType';

const ThingDescriptionTitle = (
  props: ParentProps & {
    thingType: ThingType;
    thingId: number;
    thingName: string;
  },
): JSX.Element => (
  <h3 class={styles.thing_description_title}>
    <A
      href={`${THING_STATUS_ROUTE}?thing_type=${props.thingType}&thing_id=${props.thingId}`}
      title='Открыть страницу статуса'
    >
      {props.thingId}. {props.thingName}
    </A>
  </h3>
);

export default ThingDescriptionTitle;
