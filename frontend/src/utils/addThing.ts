import { createSignal } from 'solid-js';
import axiosInstance from './axiosInstance';
import { HOME__ROUTE, ThingType } from './consts';

const addThing = async (props: {
  thing: {
    type: ThingType;
    name: string;
    userMessage: string;
    photo: string;
    location: string;
  };
}) => {
  const [data, setData] = createSignal();
  switch (props.thing.type) {
    case ThingType.lost:
      setData({
        thingType: props.thing.type,
        thingName: props.thing.name,
        userMessage: props.thing.userMessage,
        thingPhoto: props.thing.photo,
      });
    case ThingType.found:
      setData({
        thingType: props.thing.type,
        thingName: props.thing.name,
        thingLocation: props.thing.location,
        userMessage: props.thing.userMessage,
        thingPhoto: props.thing.photo,
      });
  }

  return axiosInstance
    .post(`/thing/add`, data(), {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
    .then(() => {
      window.location.href = HOME__ROUTE;
    })
    .catch((error) => console.log(error));
};

export default addThing;
