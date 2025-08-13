import { JSX, ParentProps } from 'solid-js';

import styles from './Thing.module.css';

const Thing = (props: ParentProps): JSX.Element => (
  <div class={styles.thing}>{props.children}</div>
);

export default Thing;
