import { JSX, ParentProps } from 'solid-js';

import styles from './Error.module.css';

const Error = (props: ParentProps): JSX.Element => (
  <div class={styles.error}>
    <img src='/src/assets/error.svg' />
    Ошибка!
    {props.children && ' ' + props.children}
  </div>
);

export default Error;
