import GET from "../../../utils/GET";
import { createResource } from "solid-js";

export const [lostThingsList, { refetch: syncLostThingsList }] = createResource(
  "get_things_list?type=lost",
  GET,
);

export const [foundThingsList, { refetch: syncFoundThingsList }] =
  createResource("get_things_list?type=found", GET);
