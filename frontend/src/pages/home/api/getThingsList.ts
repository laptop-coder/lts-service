import { GET } from '../../../shared/lib/utils/index';
import { createResource } from 'solid-js';

export const getThingsList = (type: string) => {
  const [thingsList, { refetch: reloadThingsList }] = createResource(
    `get_things_list?type=${type}`,
    GET,
  );
  return [thingsList, reloadThingsList];
};
