import { JSX } from 'solid-js';

import styles from './AttachFile.module.css';

const AttachFile = (
  props: JSX.InputHTMLAttributes<HTMLInputElement> & { label: string },
): JSX.Element => (
  <>
    <input
      type='file'
      accept={props.accept}
      id={props.id}
      class={styles.attach_file}
      onchange={props.onchange}
      tabindex='-1'
    />
    <label
      for={props.id}
      class={styles.attach_file_label}
      tabindex='0'
    >
      {props.label}
    </label>
  </>
);

export default AttachFile;
