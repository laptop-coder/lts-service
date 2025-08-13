import { JSX } from 'solid-js';

import styles from './Loading.module.css';

const Loading = (): JSX.Element => (
  <img
    src='/src/assets/loading.svg'
    class={styles.loading}
  />
);

export default Loading;
