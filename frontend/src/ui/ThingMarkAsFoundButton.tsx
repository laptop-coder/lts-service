import { JSX } from 'solid-js';

import { ASSETS_ROUTE } from '../utils/consts';
import thingMarkAsFound from '../utils/thingMarkAsFound';
import ThingContainerButton from './ThingContainerButton/ThingContainerButton';

const ThingMarkAsFoundButton = (props: {
  thing: {
    id: string;
    name: string;
  };
  reload?: Function;
}): JSX.Element => (
  <ThingContainerButton
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
    name='thing_mark_as_found_button'
    pathToImage={`${ASSETS_ROUTE}/yes.svg`}
    border
  />
);

export default ThingMarkAsFoundButton;
