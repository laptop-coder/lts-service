import { JSX } from 'solid-js';

import styles from './AuthInput.module.css';

const AuthInput = (
  props: JSX.InputHTMLAttributes<HTMLInputElement>,
): JSX.Element => (
  <div class={styles.auth_input_group}>
    <input
      class={styles.auth_input}
      type={props.type || 'text'}
      // Placeholder is needed here to check if the input field is not empty.
      // Using :not(:placeholder-shown)
      placeholder=''
      value={props.value}
      oninput={props.oninput}
      id={props.id}
    />
    <label
      class={styles.auth_input_label}
      for={props.id}
    >
      {props.placeholder}
    </label>
    <span class={styles.auth_input_underline} />
  </div>
);

export default AuthInput;
