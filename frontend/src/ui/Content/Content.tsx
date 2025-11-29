import { JSX, ParentProps } from 'solid-js';

import styles from './Content.module.css';

const Content = (props: ParentProps): JSX.Element => {
  return (
    <div
      class={styles.content}
      id='content'
    >
      {props.children}
    </div>
  );
};

export default Content;
