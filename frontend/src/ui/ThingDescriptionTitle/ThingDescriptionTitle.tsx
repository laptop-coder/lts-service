import { JSX, ParentProps } from 'solid-js';

import styles from './ThingDescriptionTitle.module.css';

const ThingDescriptionTitle = (
  props: ParentProps & { thingId: number; thingName: string },
): JSX.Element => (
  <h3 class={styles.thing_description_title}>
    {props.thingId}. {props.thingName}
  </h3>
);

export default ThingDescriptionTitle;
