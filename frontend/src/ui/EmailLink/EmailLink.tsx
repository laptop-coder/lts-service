import { JSX } from 'solid-js';

import styles from './EmailLink.module.css';
import type Email from '../../types/email';

const EmailLink = (props: { userEmail: Email }): JSX.Element => {
  /**
   * EmailLink
   *
   * @param userEmail<Email> - Email
   * @returns <JSX.Element> - Returns JSX element with the "mailto" link.
   */
  return (
    <a
      href={`mailto:${props.userEmail}`}
      class={styles.email_link}
      title='Написать автору объявления'
    >
      {props.userEmail}
    </a>
  );
};

export default EmailLink;
