import { JSX } from 'solid-js';

import styles from './Loading.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';

const Loading = (): JSX.Element => (
  <img
    src={`${ASSETS_ROUTE}/loading.svg`}
    class={styles.loading}
  />
);

export default Loading;
