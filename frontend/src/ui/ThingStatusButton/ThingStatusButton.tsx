import { JSX } from 'solid-js';

import styles from './ThingStatusButton.module.css';
import { ASSETS_ROUTE, USER__THING_STATUS__ROUTE } from '../../utils/consts';

const ThingStatusButton = (props: { thingId: string }): JSX.Element => (
  <button
    class={styles.thing_status_button}
    onclick={() => {
      window.location.href = `${USER__THING_STATUS__ROUTE}?thing_id=${props.thingId}`;
    }}
    title='Статус объявления'
  >
    <img src={`${ASSETS_ROUTE}/info.svg`} />
  </button>
);

export default ThingStatusButton;
