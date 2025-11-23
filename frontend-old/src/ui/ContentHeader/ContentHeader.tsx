import { JSX, ParentProps } from 'solid-js';

import styles from './ContentHeader.module.css';

const ContentHeader = (props: ParentProps): JSX.Element => (
  <h2 class={styles.content_header}>{props.children}</h2>
);

export default ContentHeader;
