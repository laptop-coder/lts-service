import { JSX, ParentProps } from 'solid-js';

import styles from './SingleList.module.css';
import List from '../List/List';

const SingleList = (props: ParentProps & { title?: string }): JSX.Element => (
  <List
    title={props.title}
    class={styles.single_list}
  >
    {props.children}
  </List>
);

export default SingleList;
