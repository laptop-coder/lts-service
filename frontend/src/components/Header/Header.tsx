import { JSX, ParentProps } from 'solid-js';

import { A } from '@solidjs/router';

import styles from './Header.module.css';
import Logo from '../../ui/Logo/Logo';
import { HOME_ROUTE } from '../../utils/consts';

const Header = (props: ParentProps): JSX.Element => {
  return (
    <header class={styles.header}>
      {/* Main module is the left part of the header*/}
      <div class={styles.main_module}>
        <Logo />
        <A href={HOME_ROUTE}>
          <h1>Сервис поиска потерянных вещей</h1>
        </A>
      </div>
      {/* Custom module is the right part of the header*/}
      <div class={styles.custom_module}>{props.children}</div>
    </header>
  );
};
export default Header;
