import { JSX } from 'solid-js';

import styles from './Error.module.css';

const Error = (): JSX.Element => (
  <div class={styles.error}>
    <img src='/src/assets/error.svg' />
    Ошибка!
  </div>
);

export default Error;
