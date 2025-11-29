import { createSignal } from 'solid-js';
import axiosInstance from './axiosInstance';
import type { LostThing, FoundThing } from '../types/thing';
import { HOME__ROUTE, ThingType } from './consts';

const addThing = async (props: {
  thingType: ThingType;
  thing: LostThing & FoundThing;
}) => {
  const [data, setData] = createSignal();
  switch (props.thingType) {
    case ThingType.lost:
      setData({
        thingType: props.thingType,
        thingName: props.thing.Name,
        userMessage: props.thing.UserMessage,
        thingPhoto: props.thing.Photo,
      });
    case ThingType.found:
      setData({
        thingType: props.thingType,
        thingName: props.thing.Name,
        thingLocation: props.thing.Location,
        userMessage: props.thing.UserMessage,
        thingPhoto: props.thing.Photo,
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
