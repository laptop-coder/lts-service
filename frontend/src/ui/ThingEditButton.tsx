import { JSX } from 'solid-js';

import { ASSETS_ROUTE, USER__THING_EDIT__ROUTE } from '../utils/consts';
import ThingContainerButton from './ThingContainerButton/ThingContainerButton';

const ThingEditButton = (props: { thing: { id: string } }): JSX.Element => (
  <ThingContainerButton
    onclick={() => {
      window.location.href = `${USER__THING_EDIT__ROUTE}?thing_id=${props.thing.id}`;
    }}
    title='Редактировать объявление'
    name='thing_edit_button'
    border
    pathToImage={`${ASSETS_ROUTE}/edit.svg`}
  />
);

export default ThingEditButton;
