import axiosInstance from './axiosInstance';
import { HOME__ROUTE, Role } from './consts';

const deleteThing = async (props: { thingId: string; role: Role }) => {
  if (props.role === Role.none) {
    return; // TODO: return error here
  }

  return axiosInstance
    .post(
      `/thing/delete/${props.role}`,
      { thingId: props.thingId },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
    .then(() => {
      window.location.href = HOME__ROUTE;
    })
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );
};

export default deleteThing;
