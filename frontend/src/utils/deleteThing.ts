import axiosInstance from './axiosInstance';
import { HOME__ROUTE, ThingType, Role } from './consts';

const deleteThing = async (props: {
  thingType: ThingType;
  thingId: string;
  role: Role;
}) => {
  if (props.role === Role.none) {
    return; // TODO: return error here
  }

  return axiosInstance
    .post(
      `/thing/delete/${props.role}`,
      { thingType: props.thingType, thingId: props.thingId },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
    .then(() => {
      window.location.href = HOME__ROUTE;
    })
    .catch((error) => console.log(error));
};

export default deleteThing;
