import { JSX } from 'solid-js';

import styles from './ThingContainerButton.module.css';

const ThingContainerButton = (props: {
  action: () => void;
  title: string;
  pathToImage: string;
}): JSX.Element => (
  <button
    class={styles.thing_container_button}
    onclick={props.action}
    title={props.title}
  >
    <img src={props.pathToImage} />
  </button>
);

export default ThingContainerButton;
