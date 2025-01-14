import type { Component } from "solid-js";
import { children } from "solid-js";
import type { DialogBox } from "../types/DialogBox";

export const DialogBox: Component = (props: DialogBox) => {
  return (
    <div class="dialog_box__wrapper">
      <div class="dialog_box">{children(() => props.children)}</div>
      <div
        class="dialog_box__background"
        onClick={props.actionToClose}
      />
    </div>
  );
};
