import { JSX, ParentProps } from 'solid-js';

import styles from './ThingDescriptionGroup.module.css';

const ThingDescriptionGroup = (props: ParentProps): JSX.Element => (
  <div class={styles.thing_description_group}>{props.children}</div>
);

export default ThingDescriptionGroup;
