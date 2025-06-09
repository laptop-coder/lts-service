import type { Component } from 'solid-js';
import type { SVGProps } from '../model/index';

export const SVG: Component<SVGProps> = (props) => {
  return (
    <svg
      class={props.class}
      xmlns='http://www.w3.org/2000/svg'
      viewBox='0 -960 960 960'
    >
      <path d={props.d} />
    </svg>
  );
};
