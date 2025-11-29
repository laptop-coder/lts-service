import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import styles from './Header.module.css';
import Logo from '../Logo/Logo';
import { Role, HOME__ROUTE, MODERATOR__HOME__ROUTE } from '../../utils/consts';
import ProfileButton from '../ProfileButton/ProfileButton';
import ThingAddButton from '../ThingAddButton/ThingAddButton';

const Header = (props: { role: Role; authorized: boolean }): JSX.Element => (
  <header class={styles.header}>
    <div class={styles.header_wrapper}>
      <Logo role={props.role} />
      <A
        href={
          props.role === Role.moderator ? MODERATOR__HOME__ROUTE : HOME__ROUTE
        }
        title='На главную'
      >
        <h1>
          <span style={{ color: 'var(--gray)' }}>LTS</span>
          <span style={{ color: 'var(--white)' }}>-сервис</span>
        </h1>
      </A>
    </div>
    <div class={styles.header_buttons}>
      {props.role === Role.user && <ThingAddButton />}
      <ProfileButton
        role={props.role}
        authorized={props.authorized}
      />
    </div>
  </header>
);

export default Header;
