import { JSX } from 'solid-js';

import styles from './HeaderLogoutButton.module.css';
import { ASSETS_ROUTE, Role } from '../../utils/consts';
import logout from '../../utils/logout';

const HeaderLogoutButton = (props: { role: Role }): JSX.Element => (
  <button
    class={styles.header_logout_button}
    title='Выйти из аккаунта'
    onclick={() => {
      if (confirm('Подтвердите выход из аккаунта')) {
        logout({ role: props.role });
      }
    }}
  >
    <img src={`${ASSETS_ROUTE}/logout.svg`} />
  </button>
);

export default HeaderLogoutButton;
