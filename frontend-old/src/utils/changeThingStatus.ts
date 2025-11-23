import ThingType from '../types/ThingType';
import axiosInstanceUnauthorized from './axiosInstanceUnauthorized';

const changeThingStatus = (props: {
  thingType: ThingType;
  thingId: number;
}) => {
  axiosInstanceUnauthorized.post(
    '/thing/change_status',
    {
      thingType: props.thingType,
      thingId: props.thingId,
    },
    {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    },
  );
};

export default changeThingStatus;
