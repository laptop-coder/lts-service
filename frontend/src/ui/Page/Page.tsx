import { JSX, ParentProps } from 'solid-js';

import styles from './Page.module.css';
import Unauthorized from '../Unauthorized/Unauthorized';
import { Role } from '../../utils/consts';

const Page = (
  props: ParentProps & { role: Role; authorized: boolean },
): JSX.Element => {
  return (
    <div class={styles.page}>
      {props.role === Role.none ||
      (props.authorized &&
        (props.role === Role.user || props.role === Role.moderator)) ? (
        props.children
      ) : (
        <Unauthorized role={props.role} />
      )}
    </div>
  );
};

export default Page;
