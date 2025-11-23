import { JSX, ParentProps } from 'solid-js';

import styles from './AuthFormAnotherAction.module.css';

const AuthFormAnotherAction = (props: ParentProps): JSX.Element => {
  return <div class={styles.auth_form_another_action}>{props.children}</div>;
};

export default AuthFormAnotherAction;
