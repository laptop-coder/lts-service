import { JSX } from 'solid-js';

import styles from './AuthForm.module.css';

const AuthForm = (
  props: JSX.FormHTMLAttributes<HTMLFormElement> & { title: string },
): JSX.Element => (
  <form
    method='post'
    onsubmit={props.onsubmit}
    class={styles.auth_form}
  >
    <h2>{props.title}</h2>
    {props.children}
  </form>
);

export default AuthForm;
