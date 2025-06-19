import type { Component } from 'solid-js';
import { children } from 'solid-js';
import type { DialogBoxProps } from '../model/index';
import { autofocus } from '@solid-primitives/autofocus';

const keyDown = (event: KeyboardEvent, actionToClose) => {
  switch (event.key) {
    case 'Escape':
      actionToClose();
      break;
  }
};

export const DialogBox: Component<DialogBoxProps> = (props) => {
  return (
    <div
      class='dialog_box__wrapper'
      tabIndex='1'
      autofocus // required for use:autofocus
      ref={autofocus}
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
