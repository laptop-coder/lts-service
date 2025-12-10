import { JSX } from 'solid-js';

import styles from './HeaderThingAddButton.module.css';
import {
  USER__THING_ADD__ROUTE,
  ASSETS_ROUTE,
  ThingType,
} from '../../utils/consts';

const HeaderThingAddButton = (props: {
  defaultThingType: ThingType;
}): JSX.Element => (
  <button
    class={styles.header_thing_add_button}
    title='Создать объявление'
    onclick={() =>
      (window.location.href =
        USER__THING_ADD__ROUTE +
        `?default_thing_type=${props.defaultThingType}`)
    }
  >
    <img src={`${ASSETS_ROUTE}/add.svg`} />
  </button>
);

export default HeaderThingAddButton;
