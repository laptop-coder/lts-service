import GET from "../../../utils/GET";

export const changeThingStatus = (id: number) => {
  GET(`change_thing_status?type=found&id=${id}`);
};
