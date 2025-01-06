import type { JSX, Component } from "solid-js";
import { children } from "solid-js";

interface DialogBoxProps {
  children: JSX.Element;
  actionToClose: func;
}

const DialogBox = (props: DialogBoxProps) => {
  return (
    <div class="dialog_box__wrapper">
      <div class="dialog_box">{children(() => props.children)}</div>
      <div
        class="dialog_box__background"
        onClick={props.actionToClose}
      ></div>
    </div>
  );
};

export default DialogBox;
