import { GET } from '../../../shared/lib/utils/index';
import { createResource } from 'solid-js';

export const getThingData = (type: string, id: number) => {
  const [thingData, { refetch: reloadThingData }] = createResource(
    `get_thing_data?type=${type}&id=${id}`,
    GET,
  );
  return [thingData, reloadThingData];
};
