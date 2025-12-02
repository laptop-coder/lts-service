import { JSX } from 'solid-js';

import styles from './HeaderLoginButton.module.css';
import {
  LOGIN_USER__ROUTE,
  LOGIN_MODERATOR__ROUTE,
  ASSETS_ROUTE,
  Role,
} from '../../utils/consts';

const HeaderLoginButton = (props: { role: Role }): JSX.Element => (
  <button
    class={styles.header_login_button}
    title='Войти в аккаунт'
    onclick={() =>
      (window.location.href = `${props.role === Role.user || props.role === Role.none ? LOGIN_USER__ROUTE : ''}${props.role === Role.moderator ? LOGIN_MODERATOR__ROUTE : ''}`)
    }
  >
    <img src={`${ASSETS_ROUTE}/login.svg`} />
  </button>
);

export default HeaderLoginButton;
