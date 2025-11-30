import { JSX, ParentProps } from 'solid-js';

import styles from './AuthFormOtherChoice.module.css';

const AuthFormOtherChoice = (props: ParentProps): JSX.Element => (
  <div class={styles.auth_form_other_choice}>{props.children}</div>
);

export default AuthFormOtherChoice;
