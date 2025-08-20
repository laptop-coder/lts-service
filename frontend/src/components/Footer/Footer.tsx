import { JSX } from 'solid-js';

import { A } from '@solidjs/router';

import styles from './Footer.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';
import { SOURCE_CODE_URL } from '../../utils/consts';

const Footer = (): JSX.Element => (
  <footer class={styles.footer}>
    <A
      href={SOURCE_CODE_URL}
      title='Исходный код'
    >
      <img src={`${ASSETS_ROUTE}/github.svg`} />
    </A>
  </footer>
);

export default Footer;
