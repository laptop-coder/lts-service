import { JSX } from 'solid-js';

import styles from './ToggleSwitch.module.css';

const ToggleSwitch = (
  props: JSX.InputHTMLAttributes<HTMLInputElement> & { sliderText?: string },
): JSX.Element => (
  <label
    for={props.id}
    class={styles.toggle_switch}
    title={props.title}
  >
    <input
      type='checkbox'
      id={props.id}
      class={styles.checkbox}
      checked={props.checked}
      oninput={props.oninput}
    />
    <span class={styles.slider_wrapper} />
    <span class={styles.slider}>{props.sliderText}</span>
  </label>
);

export default ToggleSwitch;
