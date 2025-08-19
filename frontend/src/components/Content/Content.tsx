import { JSX, ParentProps } from 'solid-js';

import styles from './Content.module.css';

const Content = (props: ParentProps & { class?: string }): JSX.Element => (
  <div class={`${styles.content} ${props.class}`}>{props.children}</div>
);

export default Content;
