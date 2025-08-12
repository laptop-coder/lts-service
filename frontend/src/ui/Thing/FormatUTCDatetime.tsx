import { JSX } from 'solid-js';

import styles from './Thing.module.css';
import type utcDatetime from '../../types/utcDatetime';
import { months } from '../../utils/consts';

const FormatUTCDatetime = (props: { datetime: utcDatetime }): JSX.Element => {
  /**
   * FormatUTCDatetime
   *
   * @param datetime<utcDatetime> - Time in the UTC format
   * @returns <JSX.Element> - Returns JSX element with the formatted datetime
   * from the UTC to the format like "01 января, 00:00". When you hover the
   * cursor over this element, the cursor style changes (to "help") and a hint
   * with a more detailed datetime format like "01.01.1970, 00:00:00" is
   * displayed.
   */
  const dt = new Date(props.datetime);
  return (
    <div
      class={styles.datetime}
      title={dt.toLocaleString()}
    >
      {dt.getDate().toString().padStart(2, '0')} {months[dt.getMonth()]},
      {dt.getHours().toString().padStart(2, '0')}:
      {dt.getMinutes().toString().padStart(2, '0')}
    </div>
  );
};

export default FormatUTCDatetime;
