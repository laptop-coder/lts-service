import { JSX } from 'solid-js';

import {
  USER__PROFILE__ROUTE,
  MODERATOR__PROFILE__ROUTE,
  ASSETS_ROUTE,
  Role,
} from '../utils/consts';
import HeaderButton from './HeaderButton/HeaderButton';

const HeaderProfileButton = (props: { role: Role }): JSX.Element => (
  <HeaderButton
    name='header_profile_button'
    title='Перейти в профиль'
    onclick={() =>
      (window.location.href = `${props.role === Role.user ? USER__PROFILE__ROUTE : ''}${props.role === Role.moderator ? MODERATOR__PROFILE__ROUTE : ''}`)
    }
    pathToImage={`${ASSETS_ROUTE}/profile.svg`}
  />
);

export default HeaderProfileButton;
