import { GET } from '@/shared/lib/utils/index';
import { createResource } from 'solid-js';

export const getThingData = (type: string, id: number) => {
  const [thingData, { refetch: reloadThingData }] = createResource(
    `thing/get_data?thing_id=${id}&thing_type=${type}`,
    GET,
  );
  return [thingData, reloadThingData];
};
