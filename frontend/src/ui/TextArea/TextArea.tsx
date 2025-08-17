import { JSX } from 'solid-js';

import styles from './TextArea.module.css';

const TextArea = (
  props: JSX.TextareaHTMLAttributes<HTMLTextAreaElement>,
): JSX.Element => (
  <textarea
    class={styles.textarea}
    placeholder={props.placeholder}
    value={props.value}
    onchange={props.onchange}
  />
);

export default TextArea;
