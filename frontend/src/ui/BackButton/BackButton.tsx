import { JSX } from 'solid-js';

import styles from './BackButton.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';

const BackButton = (): JSX.Element => (
  <button
    class={styles.back_button}
    title='Вернуться на предыдущую страницу'
    onclick={() => history.back()}
  >
    <img src={`${ASSETS_ROUTE}/arrow_back.svg`} />
  </button>
);

export default BackButton;
