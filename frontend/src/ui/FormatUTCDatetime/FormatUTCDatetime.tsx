import { JSX } from 'solid-js';

import styles from './FormatUTCDatetime.module.css';
import type UTCDatetime from '../../types/utcDatetime';
import { months } from '../../utils/consts';

const FormatUTCDatetime = (props: { datetime: UTCDatetime }): JSX.Element => {
  /**
   * FormatUTCDatetime
   *
   * @param datetime<UTCDatetime> - Time in the UTC format
   * @returns <JSX.Element> - Returns JSX element with the formatted datetime
   * from the UTC to the format like "01 января, 00:00". When you hover the
   * cursor over this element, the cursor style changes (to "help") and a hint
   * with a more detailed datetime format like "01.01.1970, 00:00:00" is
   * displayed.
   */

  /*TODO: is it normal to create new Date every time? If assign it to a
    variable, it doesn't updates when the props update.*/
  return (
    <span
      class={styles.datetime}
      title={new Date(props.datetime).toLocaleString()}
    >
      {new Date(props.datetime).getDate().toString().padStart(2, '0')}{' '}
      {months[new Date(props.datetime).getMonth()]}
      {', '}
      {new Date(props.datetime).getHours().toString().padStart(2, '0')}
      {':'}
      {new Date(props.datetime).getMinutes().toString().padStart(2, '0')}
    </span>
  );
};

export default FormatUTCDatetime;
