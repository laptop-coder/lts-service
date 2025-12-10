import { JSX } from 'solid-js';

import styles from './ThingEditButton.module.css';
import { ASSETS_ROUTE, USER__THING_EDIT__ROUTE } from '../../utils/consts';

const ThingEditButton = (props: { thing: { id: string } }): JSX.Element => (
  <button
    class={styles.thing_edit_button}
    onclick={() => {
      window.location.href = `${USER__THING_EDIT__ROUTE}?thing_id=${props.thing.id}`;
    }}
    title='Редактировать объявление'
  >
    <img src={`${ASSETS_ROUTE}/edit.svg`} />
  </button>
);

export default ThingEditButton;
