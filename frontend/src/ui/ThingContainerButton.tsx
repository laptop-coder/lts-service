import { JSX } from 'solid-js';

import styles from './ThingContainerButton.module.css';

const ThingContainerButton = (props: {
  title: string;
  pathToImage: string;
  type?: 'button' | 'submit' | 'reset' | 'menu' | undefined;
  onclick?: () => void;
  name: string;
  border?: boolean;
}): JSX.Element => (
  <button
    name={props.name}
    class={`${styles.thing_container_button} ${props.border ? styles.border : ''}`}
    onclick={props.onclick}
    title={props.title}
    type={props.type || 'button'}
  >
    <img src={props.pathToImage} />
  </button>
);

export default ThingContainerButton;
