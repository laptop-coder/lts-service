import { JSX, ParentProps } from 'solid-js';

import styles from './FormTitle.module.css';

const FormTitle = (props: ParentProps): JSX.Element => (
  <h2 class={styles.form_title}>{props.children}</h2>
);

export default FormTitle;
