import { JSX } from 'solid-js';

import styles from './SubmitButton.module.css';

const SubmitButton = (
  props: JSX.ButtonHTMLAttributes<HTMLButtonElement>,
): JSX.Element => (
  <button
    class={styles.submit_button}
    type='submit'
  >
    {props.children}
  </button>
);

export default SubmitButton;
