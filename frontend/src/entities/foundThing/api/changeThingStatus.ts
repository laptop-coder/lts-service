import { GET } from "../../../shared/lib/utils/index";

export const changeThingStatus = (id: number) => {
  return new Promise(async (res, ref) => {
    await GET(`change_thing_status?type=found&id=${id}`);
    res();
  });
};
