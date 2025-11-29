import { JSX, ParentProps } from 'solid-js';

import styles from './TabsToggleWrapper.module.css';

const TabsToggleWrapper = (
  props: ParentProps & { fullscreen?: boolean },
): JSX.Element => {
  return (
    <div
      class={
        props.fullscreen
          ? styles.tabs_toggle_wrapper_fullscreen
          : styles.tabs_toggle_wrapper
      }
    >
      {props.children}
    </div>
  );
};

export default TabsToggleWrapper;
