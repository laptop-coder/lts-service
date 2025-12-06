import axiosInstance from './axiosInstance';
import { HOME__ROUTE, ThingType } from './consts';

const editThing = async (props: {
  thing: {
    id: string;
    newType: ThingType;
    newName: string;
    newUserMessage: string;
    newPhoto: string;
  };
}) => {
  const data = {
    thingId: props.thing.id,
    newThingType: props.thing.newType,
    newThingName: props.thing.newName,
    newUserMessage: props.thing.newUserMessage,
    newThingPhoto: props.thing.newPhoto,
  };

  return axiosInstance
    .post(`/thing/edit`, data, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
    .then(() => {
      window.location.href = HOME__ROUTE;
    })
    .catch((error) =>
      // console.log(error) // TODO: think about it, how to make right
      console.log('error'),
    );
};

export default editThing;
