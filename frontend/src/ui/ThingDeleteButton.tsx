import { JSX } from 'solid-js';

import { ASSETS_ROUTE, Role } from '../utils/consts';
import deleteThing from '../utils/deleteThing';
import ThingContainerButton from './ThingContainerButton/ThingContainerButton';

const ThingDeleteButton = (props: {
  thing: { id: string; name: string };
  role: Role;
}): JSX.Element => (
  <ThingContainerButton
    onclick={() => {
      if (
        confirm(
          `Подтвердите удаление объявления "${props.thing.name}". Это действие необратимо`,
        )
      ) {
        deleteThing({ thing: { id: props.thing.id }, role: props.role });
      }
    }}
    title='Удалить объявление'
    pathToImage={`${ASSETS_ROUTE}/delete.svg`}
    name='thing_container_button'
    border
  />
);

export default ThingDeleteButton;
