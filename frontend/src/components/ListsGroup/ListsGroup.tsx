import { JSX, ParentProps } from 'solid-js';

import styles from './ListsGroup.module.css';

const ListsGroup = (props: ParentProps): JSX.Element => (
  <div class={styles.lists_group}>{props.children}</div>
);

export default ListsGroup;
