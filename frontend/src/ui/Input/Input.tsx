import { JSX } from 'solid-js';

import styles from './Input.module.css';

const Input = (
  props: JSX.InputHTMLAttributes<HTMLInputElement>,
): JSX.Element => (
  <input
    class={styles.input}
    type={props.type || 'text'}
    placeholder={props.placeholder}
    value={props.value}
    onchange={props.onchange}
  />
);

export default Input;
