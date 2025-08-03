import { POST } from '@/shared/lib/utils/index';

export const changeThingStatus = (thingId: number, thingType: string) => {
  return new Promise((resolve, reject) => {
    POST(`thing/change_status`, {
      thing_id: thingId,
      thing_type: thingType,
    });
  });
};
