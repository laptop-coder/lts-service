import { JSX } from 'solid-js';

import styles from './Unauthorized.module.css';
import {
  LOGIN_USER__ROUTE,
  REGISTER_USER__ROUTE,
  LOGIN_MODERATOR__ROUTE,
  REGISTER_MODERATOR__ROUTE,
  Role,
} from '../../utils/consts';

import { A } from '@solidjs/router';

const Unauthorized = (props: { role: Role }): JSX.Element => {
  return (
    <div class={styles.unauthorized}>
      {props.role === Role.user && (
        <>
          Для доступа к этой странице необходимо зарегистрироваться.{' '}
          <A href={LOGIN_USER__ROUTE}>Войдите</A> в аккаунт или{' '}
          <A href={REGISTER_USER__ROUTE}>создайте</A> новый.
        </>
      )}
      {props.role === Role.moderator && (
        <>
          Для доступа к этой странице необходим аккаунт модератора.{' '}
          <A href={LOGIN_MODERATOR__ROUTE}>Войдите</A> в аккаунт или{' '}
          <A href={REGISTER_MODERATOR__ROUTE}>создайте</A> новый. Обратите
          внимание, что в системе может быть создан только один аккаунт
          модератора.
        </>
      )}
    </div>
  );
};

export default Unauthorized;
