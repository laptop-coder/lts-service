import { JSX, ParentProps } from 'solid-js';

import styles from './ThingDescriptionItem.module.css';

const ThingDescriptionItem = (props: ParentProps): JSX.Element => (
  <div class={styles.thing_description_item}>{props.children}</div>
);

export default ThingDescriptionItem;
