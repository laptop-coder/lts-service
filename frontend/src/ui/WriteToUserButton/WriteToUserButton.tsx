import { JSX } from 'solid-js';

import styles from './WriteToUserButton.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';

const WriteToUserButton = (props: { email: string }): JSX.Element => (
  <button
    class={styles.write_to_user_button}
    type='button'
    name='write_to_user_button'
    onclick={() => (window.location.href = 'mailto:' + props.email)}
  >
    <img src={`${ASSETS_ROUTE}/message.svg`} />
  </button>
);

export default WriteToUserButton;
