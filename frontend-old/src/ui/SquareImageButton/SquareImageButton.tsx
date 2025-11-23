import { JSX } from 'solid-js';

import styles from './SquareImageButton.module.css';

const SquareImageButton = (
  props: JSX.ButtonHTMLAttributes<HTMLButtonElement>,
): JSX.Element => (
  <button
    class={styles.square_image_button}
    onclick={props.onclick}
    type='button'
  >
    {props.children}
  </button>
);

export default SquareImageButton;
