import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import { HOME_ROUTE, MODERATOR_ROUTE } from '../../utils/consts';
import { ASSETS_ROUTE } from '../../utils/consts';
import styles from './Logo.module.css';

const Logo = (props: { moderator?: boolean }): JSX.Element => {
  return (
    <A
      class={styles.logo}
      href={props.moderator ? MODERATOR_ROUTE : HOME_ROUTE}
      title='На главную'
    >
      <img src={`${ASSETS_ROUTE}/logo.svg`} />
    </A>
  );
};

export default Logo;
