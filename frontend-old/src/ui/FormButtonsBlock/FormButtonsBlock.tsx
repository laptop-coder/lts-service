import { JSX, ParentProps } from 'solid-js';

import styles from './FormButtonsBlock.module.css';

const FormButtonsBlock = (props: ParentProps): JSX.Element => (
  <div class={styles.form_buttons_block}>{props.children}</div>
);

export default FormButtonsBlock;
