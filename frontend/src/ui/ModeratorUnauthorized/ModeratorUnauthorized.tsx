import { JSX } from 'solid-js';

import styles from './ModeratorUnauthorized.module.css';
import { A } from '@solidjs/router';

const ModeratorUnauthorized = (): JSX.Element => (
  <div class={styles.moderator_unauthorized}>
    Для получения доступа к этой странице{' '}
    <A href='/moderator/login'>войдите в аккаунт</A> модератора или{' '}
    <A href='/moderator/register'>создайте новый</A>
  </div>
);

export default ModeratorUnauthorized;
