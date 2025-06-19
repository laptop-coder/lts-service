import { GET } from '../../../shared/lib/utils/index';

export const changeThingStatus = (id: number) => {
  return new Promise((resolve, reject) => {
    GET(`change_thing_status?type=lost&id=${id}`)
      .then((response) => resolve(response))
      .catch((error) => reject(error));
  });
};
