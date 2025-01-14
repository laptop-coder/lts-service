import { GET } from "../../../shared/lib/utils/index";
import { createResource } from "solid-js";

export const [lostThingsList, { refetch: syncLostThingsList }] = createResource(
  "get_things_list?type=lost",
  GET,
);

export const [foundThingsList, { refetch: syncFoundThingsList }] =
  createResource("get_things_list?type=found", GET);
