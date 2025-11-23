import { JSX } from 'solid-js';

import styles from './Select.module.css';

const Select = (
  props: JSX.SelectHTMLAttributes<HTMLSelectElement> & { label: string },
): JSX.Element => (
  <div class={styles.select_wrapper}>
    <label
      class={styles.select_label}
      for={props.name}
    >
      {props.label}
    </label>
    <select
      class={styles.select}
      id={props.id}
      oninput={props.oninput}
    >
      {props.children}
    </select>
  </div>
);

export default Select;
