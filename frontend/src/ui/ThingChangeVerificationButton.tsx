import { JSX } from 'solid-js';

import { ASSETS_ROUTE, VerificationAction } from '../utils/consts';
import thingChangeVerification from '../utils/thingChangeVerification';
import ThingContainerButton from './ThingContainerButton/ThingContainerButton';

const ThingChangeVerificationButton = (props: {
  thing: { id: string };
  action: VerificationAction;
  reload?: Function;
}): JSX.Element => {
  return (
    <ThingContainerButton
      green={props.action === VerificationAction.approve}
      red={props.action === VerificationAction.reject}
      onclick={() => {
        thingChangeVerification({
          thing: { id: props.thing.id },
          action: props.action,
        });
        if (props.reload) {
          props.reload();
        }
      }}
      title={
        props.action === VerificationAction.approve
          ? 'Одобрить объявление'
          : 'Отклонить объявление'
      }
      name='thing_container_button'
      pathToImage={`${ASSETS_ROUTE}/${props.action === VerificationAction.approve ? 'yes' : ''}${props.action === VerificationAction.reject ? 'no' : ''}.svg`}
    />
  );
};

export default ThingChangeVerificationButton;
