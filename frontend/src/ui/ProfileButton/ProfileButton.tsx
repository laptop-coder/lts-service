import { JSX } from 'solid-js';

import styles from './ProfileButton.module.css';
import {
  USER__PROFILE__ROUTE,
  MODERATOR__PROFILE__ROUTE,
  LOGIN_USER__ROUTE,
  LOGIN_MODERATOR__ROUTE,
  ASSETS_ROUTE,
  Role,
} from '../../utils/consts';

const ProfileButton = (props: {
  role: Role;
  authorized: boolean;
}): JSX.Element => (
  <>
    {props.authorized === true ? (
      <button
        class={styles.profile_button}
        title='Перейти в профиль'
        onclick={() =>
          (window.location.href = `${props.role === Role.user ? USER__PROFILE__ROUTE : ''}${props.role === Role.moderator ? MODERATOR__PROFILE__ROUTE : ''}`)
        }
      >
        <img src={`${ASSETS_ROUTE}/profile.svg`} />
      </button>
    ) : (
      <button
        class={styles.login_button}
        title='Войти в аккаунт'
        onclick={() =>
          (window.location.href = `${props.role === Role.user || props.role === Role.none ? LOGIN_USER__ROUTE : ''}${props.role === Role.moderator ? LOGIN_MODERATOR__ROUTE : ''}`)
        }
      >
        <img src={`${ASSETS_ROUTE}/login.svg`} />
      </button>
    )}
  </>
);

export default ProfileButton;
