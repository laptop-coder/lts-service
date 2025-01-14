import { GET } from "../../../shared/lib/utils/index";

export const changeThingStatus = (id: number) => {
  GET(`change_thing_status?type=found&id=${id}`);
};
