import { JSX } from 'solid-js';

import {
  LOGIN_USER__ROUTE,
  LOGIN_MODERATOR__ROUTE,
  ASSETS_ROUTE,
  Role,
} from '../utils/consts';
import HeaderButton from './HeaderButton/HeaderButton';

const HeaderLoginButton = (props: { role: Role }): JSX.Element => (
  <HeaderButton
    name='header_login_button'
    title='Войти в аккаунт'
    onclick={() =>
      (window.location.href = `${props.role === Role.user || props.role === Role.none ? LOGIN_USER__ROUTE : ''}${props.role === Role.moderator ? LOGIN_MODERATOR__ROUTE : ''}`)
    }
    pathToImage={`${ASSETS_ROUTE}/login.svg`}
  />
);

export default HeaderLoginButton;
