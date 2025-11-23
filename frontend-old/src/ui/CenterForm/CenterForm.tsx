import { JSX } from 'solid-js';

import styles from './CenterForm.module.css';

const CenterForm = (
  props: JSX.FormHTMLAttributes<HTMLFormElement> & { header: string },
): JSX.Element => (
  <div class={styles.center_form_wrapper}>
    <form
      class={styles.center_form}
      method={props.method}
      onsubmit={props.onsubmit}
    >
      <h2 class={styles.center_form_header}>{props.header}</h2>
      {props.children}
    </form>
  </div>
);

export default CenterForm;
