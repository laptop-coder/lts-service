import { JSX } from 'solid-js';

import styles from './ThingAddButton.module.css';
import { USER__THING_ADD__ROUTE, ASSETS_ROUTE } from '../../utils/consts';

const ThingAddButton = (): JSX.Element => (
  <button
    class={styles.thing_add_button}
    title='Создать объявление'
    onclick={() => (window.location.href = USER__THING_ADD__ROUTE)}
  >
    <img src={`${ASSETS_ROUTE}/add.svg`} />
  </button>
);

export default ThingAddButton;
