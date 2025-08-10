import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import styles from './Header.module.css';
import Logo from '../../ui/Logo/Logo';
import { HOME_ROUTE } from '../../utils/consts';

const Header = (): JSX.Element => {
  return (
    <header class={styles.header}>
      <Logo />
      <A href={HOME_ROUTE}>
        <h1>Сервис поиска потерянных вещей</h1>
      </A>
    </header>
  );
};
export default Header;
