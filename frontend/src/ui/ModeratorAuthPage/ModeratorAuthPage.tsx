import { JSX, ParentProps } from 'solid-js';

import styles from './ModeratorAuthPage.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';

const ModeratorAuthPage = (props: ParentProps): JSX.Element => {
  return (
    <div
      class={styles.moderator_auth_page}
      style={{
        'background-image': `url(${ASSETS_ROUTE}/auth_background.jpeg)`,
        'background-position': 'center',
        'background-size': 'cover',
      }}
    >
      {props.children}
    </div>
  );
};

export default ModeratorAuthPage;
