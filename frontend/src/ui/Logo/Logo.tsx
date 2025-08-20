import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import { HOME_ROUTE } from '../../utils/consts';
import { ASSETS_ROUTE } from '../../utils/consts';
import styles from './Logo.module.css';

const Logo = (): JSX.Element => {
  return (
    <A
      class={styles.logo}
      href={HOME_ROUTE}
    >
      <img src={`${ASSETS_ROUTE}/logo.svg`} />
    </A>
  );
};

export default Logo;
