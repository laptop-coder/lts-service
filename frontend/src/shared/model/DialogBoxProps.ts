import type { JSX } from 'solid-js';
export interface DialogBoxProps {
  children: JSX.Element;
  actionToClose: Function;
}
