import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import {
  HOME__ROUTE,
  MODERATOR__HOME__ROUTE,
  ASSETS_ROUTE,
  Role,
} from '../../utils/consts';
import styles from './Logo.module.css';

const Logo = (props: { role: Role }): JSX.Element => {
  return (
    <A
      class={styles.logo}
      href={
        props.role === Role.moderator ? MODERATOR__HOME__ROUTE : HOME__ROUTE
      }
      title='На главную'
    >
      <img
        class={styles.img}
        src={`${ASSETS_ROUTE}/logo.svg`}
      />
    </A>
  );
};

export default Logo;
