import { JSX } from 'solid-js';

import styles from './AuthButton.module.css';

const AuthButton = (
  props: JSX.ButtonHTMLAttributes<HTMLButtonElement>,
): JSX.Element => (
  <button
    class={styles.auth_button}
    type='submit'
  >
    {props.children}
  </button>
);

export default AuthButton;
