import { JSX } from 'solid-js';

import styles from './HeaderButton.module.css';

const HeaderButton = (props: {
  title: string;
  pathToImage: string;
  onclick?: () => void;
  name: string;
}): JSX.Element => (
  <button
    name={props.name}
    class={styles.header_button}
    onclick={props.onclick}
    title={props.title}
    type='button'
  >
    <img src={props.pathToImage} />
  </button>
);

export default HeaderButton;
