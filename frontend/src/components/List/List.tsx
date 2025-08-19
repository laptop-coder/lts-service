import { JSX, ParentProps } from 'solid-js';

import styles from './List.module.css';

const List = (
  props: ParentProps & { title?: string; class?: string },
): JSX.Element => (
  <div class={`${styles.list} ${props.class}`}>
    <h2>{props.title}</h2>
    <div class={styles.content}>{props.children}</div>
  </div>
);

export default List;
