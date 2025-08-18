import { JSX, Show } from 'solid-js';

import styles from './ThingPhoto.module.css';

const ThingPhoto = (
  props: JSX.ImgHTMLAttributes<HTMLImageElement>,
): JSX.Element => (
  <Show when={props.src}>
    <img
      src={props.src}
      title={props.title}
      class={styles.thing_photo}
    />
  </Show>
);

export default ThingPhoto;
