import type { Component } from 'solid-js';
import { children } from 'solid-js';
import type { DialogBox } from '../model/DialogBox';
import { autofocus } from '@solid-primitives/autofocus';

const keyDown = (event, actionToClose) => {
  switch (event.key) {
    case 'Escape':
      actionToClose();
      break;
  }
};

export const DialogBox: Component = (props: DialogBox) => {
  return (
    <div
      class='dialog_box__wrapper'
      tabindex='1'
      autofocus
      use:autofocus
      onKeyDown={(event) => keyDown(event, props.actionToClose)}
    >
      <div class='dialog_box'>{children(() => props.children)}</div>
      <div
        class='dialog_box__background'
        onClick={props.actionToClose}
      />
    </div>
  );
};
