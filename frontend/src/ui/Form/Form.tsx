import { JSX, ParentProps } from 'solid-js';

import styles from './Form.module.css';

const Form = (
  props: JSX.FormHTMLAttributes<HTMLFormElement> & ParentProps,
): JSX.Element => (
  <form
    class={styles.form}
    onsubmit={props.onsubmit}
  >
    {props.children}
  </form>
);

export default Form;
