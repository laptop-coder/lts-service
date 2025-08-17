import { JSX, Switch, Match } from 'solid-js';

import styles from './ThingStatus.module.css';
import TextAnimation from '../TextAnimation/TextAnimation';
import Error from '../Error/Error';

const ThingStatus = (props: {
  verified: number;
  status: number;
}): JSX.Element => (
  <Switch fallback={<Error />}>
    <Match when={props.status === 0 && props.verified === 0}>
      <TextAnimation>
        Объявление на{' '}
        <span class={styles.on_moderation_text}>
          <b>модерации</b>
        </span>
        , <b>не опубликовано</b>
      </TextAnimation>
    </Match>
    <Match when={props.status === 0 && props.verified === 1}>
      <TextAnimation>
        Объявление <b>опубликовано</b>, вещь{' '}
        <span class={styles.thing_not_found_text}>
          <b>не найдена</b>
        </span>
      </TextAnimation>
    </Match>
    <Match when={props.status === 1 && props.verified === 1}>
      <TextAnimation>
        Вещь{' '}
        <span class={styles.thing_found_text}>
          <b>найдена</b>
        </span>
        , объявление <b>снято с публикации</b>
      </TextAnimation>
    </Match>
  </Switch>
);

export default ThingStatus;
