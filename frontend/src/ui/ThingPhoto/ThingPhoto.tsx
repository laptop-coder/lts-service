import { Show, JSX, ParentProps } from 'solid-js';

import styles from './ThingPhoto.module.css';

const ThingPhoto = (
  props: JSX.InputHTMLAttributes<HTMLInputElement>,
): JSX.Element => (
  <Show when={props.src}>
    <img
      src={props.src}
      class={styles.thing_photo}
    />
  </Show>
);

export default ThingPhoto;
