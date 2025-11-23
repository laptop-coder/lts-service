import { JSX, ParentProps } from 'solid-js';

import styles from './Page.module.css';

const Page = (props: ParentProps): JSX.Element => {
  return <div class={styles.page}>{props.children}</div>;
};

export default Page;
