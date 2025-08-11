import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import { HOME_ROUTE } from '../../utils/consts';
import styles from './Logo.module.css';

const Logo = (): JSX.Element => {
  return (
    <A
      class={styles.logo}
      href={HOME_ROUTE}
    >
      <svg>
        <image href='/src/assets/logo.svg'></image>
      </svg>
    </A>
  );
};

export default Logo;
