import { JSX } from 'solid-js';

import styles from './UserProfileContent.module.css';

const UserProfileContent = (props: {
  username: string;
  email: string;
}): JSX.Element => (
  <div class={styles.user_profile_content}>
    <span>Имя пользователя: {props.username}</span>
    <span>Email: {props.email}</span>
  </div>
);

export default UserProfileContent;
