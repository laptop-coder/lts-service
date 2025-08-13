import { JSX } from 'solid-js';

import styles from './EmailLink.module.css';
import type email from '../../types/email';

const EmailLink = (props: { userEmail: email }): JSX.Element => {
  /**
   * EmailLink
   *
   * @param userEmail<email> - Email
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
