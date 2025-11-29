import { JSX, ParentProps } from 'solid-js';

import styles from './ThingContainerItem.module.css';

const ThingContainerItem = (
  props: ParentProps & { pathToImage: string; title: string },
): JSX.Element => (
  <div
    class={styles.thing_container_item}
    title={props.title}
  >
    <img
      src={props.pathToImage}
      class={styles.img}
    />
    {props.children}
  </div>
);

export default ThingContainerItem;
