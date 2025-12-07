import axiosInstance from '../utils/axiosInstance';
import { VerificationAction } from '../utils/consts';

const thingChangeVerification = (props: {
  thingId: string;
  action: VerificationAction;
}) =>
  axiosInstance
    .post(
      `/thing/change_verification`,
      {
        thingId: props.thingId,
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
      // console.log(error) // think about it, how to make right
      console.log('error'),
    );

export default thingChangeVerification;
