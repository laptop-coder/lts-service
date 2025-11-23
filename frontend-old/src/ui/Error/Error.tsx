import { JSX, ParentProps } from 'solid-js';

import styles from './Error.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';

const Error = (props: ParentProps): JSX.Element => (
  <div class={styles.error}>
    <img
      src={`${ASSETS_ROUTE}/error.svg`}
      class={styles.img}
    />
    Ошибка!
    {props.children && ' ' + props.children}
  </div>
);

export default Error;
