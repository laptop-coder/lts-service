import { JSX } from 'solid-js';

import styles from './ResetButton.module.css';

const ResetButton = (
  props: JSX.ButtonHTMLAttributes<HTMLButtonElement>,
): JSX.Element => (
  <button
    class={styles.reset_button}
    onclick={props.onclick}
    type='reset'
  >
    {props.children}
  </button>
);

export default ResetButton;
