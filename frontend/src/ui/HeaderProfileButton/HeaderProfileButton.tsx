import { JSX } from 'solid-js';

import styles from './HeaderProfileButton.module.css';
import {
  USER__PROFILE__ROUTE,
  MODERATOR__PROFILE__ROUTE,
  ASSETS_ROUTE,
  Role,
} from '../../utils/consts';

const HeaderProfileButton = (props: { role: Role }): JSX.Element => (
  <button
    class={styles.profile_button}
    title='Перейти в профиль'
    onclick={() =>
      (window.location.href = `${props.role === Role.user ? USER__PROFILE__ROUTE : ''}${props.role === Role.moderator ? MODERATOR__PROFILE__ROUTE : ''}`)
    }
  >
    <img src={`${ASSETS_ROUTE}/profile.svg`} />
  </button>
);

export default HeaderProfileButton;
