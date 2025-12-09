import { JSX } from 'solid-js';

import styles from './ThingMarkAsFoundButton.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';
import thingMarkAsFound from '../../utils/thingMarkAsFound';

const ThingMarkAsFoundButton = (props: {
  thing: {
    id: string;
    name: string;
  };
  reload?: Function;
}): JSX.Element => {
  return (
    <button
      class={styles.thing_mark_as_found_button}
      onclick={() => {
        if (
          confirm(
            `Подтвердите изменение статуса объявления "${props.thing.name}". После подтверждения вещь будет считаться найденной. Это действие необратимо.`,
          )
        ) {
          thingMarkAsFound({ thing: { id: props.thing.id } });
          if (props.reload) {
            props.reload();
          }
        }
      }}
      title='Отметить вещь как найденную'
    >
      <img src={`${ASSETS_ROUTE}/yes.svg`} />
    </button>
  );
};

export default ThingMarkAsFoundButton;
