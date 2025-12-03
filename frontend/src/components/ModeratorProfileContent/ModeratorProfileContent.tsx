import { JSX } from 'solid-js';

import styles from './ModeratorProfileContent.module.css';

const ModeratorProfileContent = (props: { username: string }): JSX.Element => (
  <div class={styles.moderator_profile_content}>
    <span>Имя пользователя (модератор): {props.username}</span>
  </div>
);

export default ModeratorProfileContent;
