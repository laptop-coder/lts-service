import { JSX } from 'solid-js';

import styles from './ThingChangeVerificationButton.module.css';
import { ASSETS_ROUTE } from '../../utils/consts';
import thingChangeVerification from '../../utils/thingChangeVerification';
import { VerificationAction } from '../../utils/consts';

const ThingChangeVerificationButton = (props: {
  thingId: string;
  action: VerificationAction;
}): JSX.Element => (
  <button
    class={`${props.action === VerificationAction.approve ? styles.thing_change_verification_button_approve : ''} ${props.action === VerificationAction.reject ? styles.thing_change_verification_button_reject : ''}`}
    onclick={() => {
      thingChangeVerification({ thingId: props.thingId, action: props.action });
    }}
    title={
      props.action === VerificationAction.approve
        ? 'Одобрить объявление'
        : 'Отклонить объявление'
    }
  >
    <img
      src={`${ASSETS_ROUTE}/${props.action === VerificationAction.approve ? 'yes' : ''}${props.action === VerificationAction.reject ? 'no' : ''}.svg`}
    />
  </button>
);

export default ThingChangeVerificationButton;
