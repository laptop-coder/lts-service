import { JSX, ParentProps } from 'solid-js';

import styles from './FormButtonsGroup.module.css';

const FormButtonsGroup = (props: ParentProps): JSX.Element => (
  <div class={styles.form_buttons_group}>{props.children}</div>
);

export default FormButtonsGroup;
