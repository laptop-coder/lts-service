import ThingType from '../types/ThingType';
import axiosInstanceUnauthorized from './axiosInstanceUnauthorized';

type ThingVerificationAction = 'accept' | 'reject';

const verifyThing = (props: {
  thingType: ThingType;
  thingId: number;
  action: ThingVerificationAction;
}) => {
  axiosInstanceUnauthorized.post(
    '/thing/verify',
    {
      thingType: props.thingType,
      thingId: props.thingId,
      action: props.action,
    },
    {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    },
  );
};

export default verifyThing;
