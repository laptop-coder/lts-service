import axiosInstance from '../utils/axiosInstance';
import { VerificationAction } from '../utils/consts';

const thingChangeVerification = (props: {
  thing: { id: string };
  action: VerificationAction;
}) =>
  axiosInstance
    .post(
      `/thing/change_verification`,
      {
        thingId: props.thing.id,
        action: props.action,
      },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
    .then(() => {
      // window.location.href = HOME__ROUTE;
    })
    .catch((error) =>
      // console.log(error) // TODO: think about it, how to make right
      console.log('error'),
    );

export default thingChangeVerification;
