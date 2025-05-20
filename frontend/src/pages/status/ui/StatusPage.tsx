import type { Component } from "solid-js";
import { useParams } from "@solidjs/router";
import "../../../app/styles.css";

export const StatusPage: Component = () => {
  const params = useParams();
  console.log(params.type, params.id);
  return <div class="page"></div>;
};
