import { GET } from '../../../shared/lib/utils/index';

export const changeThingStatus = (id: number) => {
  return new Promise(async (resolve, reject) => {
    await GET(`change_thing_status?type=found&id=${id}`);
    resolve();
  });
};
