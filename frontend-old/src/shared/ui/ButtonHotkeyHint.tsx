import type { ButtonHotkeyHintProps } from '../model/ButtonHotkeyHintProps';
import type { Component } from 'solid-js';

export const ButtonHotkeyHint: Component<ButtonHotkeyHintProps> = (props) => {
  return (
    <div class={`button__hotkey_hint ${props.place} ${props.side}`}>
      {props.hotkey}
    </div>
  );
};
