import { JSX } from 'solid-js';

import { ASSETS_ROUTE, Role } from '../utils/consts';
import logout from '../utils/logout';
import HeaderButton from './HeaderButton/HeaderButton';

const HeaderLogoutButton = (props: { role: Role }): JSX.Element => (
  <HeaderButton
    name='header_logout_button'
    title='Выйти из аккаунта'
    onclick={() => {
      if (confirm('Подтвердите выход из аккаунта')) {
        logout({ role: props.role });
      }
    }}
    pathToImage={`${ASSETS_ROUTE}/logout.svg`}
  />
);

export default HeaderLogoutButton;
