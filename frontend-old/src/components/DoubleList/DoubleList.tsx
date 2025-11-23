import { JSX, ParentProps } from 'solid-js';

import styles from './DoubleList.module.css';
import List from '../List/List';

const DoubleList = (props: ParentProps & { title?: string }): JSX.Element => (
  <List
    title={props.title}
    class={styles.double_list}
  >
    {props.children}
  </List>
);

export default DoubleList;
