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
import logout from '../../utils/logout';

const ProfileButton = (props: {
  role: Role;
  authorized: boolean;
  showLogout?: boolean;
}): JSX.Element => (
  //TODO: refactor (rewrite)
  <>
    {props.authorized === true ? (
      props.showLogout ? (
        <button
          class={styles.profile_button}
          title='Выйти из аккаунта'
          onclick={() => {
            if (confirm('Подтвердите выход из аккаунта')) {
              logout({ role: props.role });
            }
          }}
        >
          <img src={`${ASSETS_ROUTE}/logout.svg`} />
        </button>
      ) : (
        <button
          class={styles.profile_button}
          title='Перейти в профиль'
          onclick={() =>
            (window.location.href = `${props.role === Role.user ? USER__PROFILE__ROUTE : ''}${props.role === Role.moderator ? MODERATOR__PROFILE__ROUTE : ''}`)
          }
        >
          <img src={`${ASSETS_ROUTE}/profile.svg`} />
        </button>
      )
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
